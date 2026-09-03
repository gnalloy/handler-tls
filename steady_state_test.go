package tls

import (
	cryptotls "crypto/tls"
	"testing"

	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
)

func TestFinishHandshakeEnablesNonblockingInput(t *testing.T) {
	handler := newHandler(ModeServer, Config{StartTLS: true})
	handler.raw = newMemoryConn(handler.bytePool, handler.notifyDrain)
	if err := handler.finishHandshake(); err != nil {
		t.Fatal(err)
	}
	if !handler.raw.nonblocking.Load() {
		t.Fatal("nonblocking input was not enabled")
	}
}

func TestHandlerWritesApplicationDataSynchronouslyAfterHandshake(t *testing.T) {
	conn := &steadyStateConn{}
	handler, ch := newSteadyStateHandler(t, conn, nil)

	out := buffer.NewHeapBuffer(4)
	_, _ = out.WriteBytes([]byte("ping"))
	if err := ch.WriteAndFlush(out); err != nil {
		t.Fatal(err)
	}
	if got := string(conn.written); got != "ping" {
		t.Fatalf("written=%q, want ping", got)
	}
	if handler.pendingApplication.len() != 0 {
		t.Fatalf("pending application buffers=%d, want 0", handler.pendingApplication.len())
	}
}

func TestHandlerWriteAndFlushCombinesFinalCiphertextWrite(t *testing.T) {
	sink := &combinedCipherSink{}
	ch := channel.NewLocalChannel(1, buffer.NewHeapAllocator(), sink)
	handler := newHandler(ModeServer, Config{StartTLS: true})
	if err := ch.Pipeline().AddLast("tls", handler); err != nil {
		t.Fatal(err)
	}
	handler.raw = newMemoryConn(handler.bytePool, handler.notifyDrain)
	handler.raw.setNonblocking()
	handler.conn = &ciphertextConn{output: handler.raw, split: true}
	handler.startOnce.Do(func() {})
	handler.pendingApplication.markReady()
	handler.started.Store(true)
	handler.handshake = true
	handler.handshakeComplete.Store(true)
	handler.active = true
	handler.activated.Store(true)

	out := buffer.NewHeapBuffer(4)
	_, _ = out.WriteBytes([]byte("ping"))
	if err := ch.WriteAndFlush(out); err != nil {
		t.Fatal(err)
	}
	if sink.combinedWrites != 1 || sink.writes != 2 || sink.flushes != 0 {
		t.Fatalf("combined writes=%d writes=%d flushes=%d, want 1/2/0", sink.combinedWrites, sink.writes, sink.flushes)
	}
	if got := string(sink.data); got != "ping" {
		t.Fatalf("ciphertext=%q, want ping", got)
	}
}

func TestHandlerReadsApplicationDataSynchronouslyAfterHandshake(t *testing.T) {
	conn := &steadyStateConn{readData: []byte("plain")}
	recorder := &plainRecorder{}
	handler, ch := newSteadyStateHandler(t, conn, recorder)
	if !handler.handshakeComplete.Load() || handler.conn != conn {
		t.Fatalf("steady state was not initialized")
	}

	in := buffer.NewHeapBuffer(1)
	_, _ = in.WriteBytes([]byte{1})
	ch.Pipeline().FireChannelRead(in)

	if got := recorder.String(); got != "plain" {
		t.Fatalf("plaintext=%q reads=%d, want plain", got, conn.reads)
	}
	if conn.reads != 2 {
		t.Fatalf("reads=%d, want one data read and one would-block read", conn.reads)
	}
}

func newSteadyStateHandler(t *testing.T, conn Conn, inbound channel.Handler) (*Handler, channel.Channel) {
	t.Helper()
	sink := &pipeSink{}
	ch := channel.NewLocalChannel(1, buffer.NewHeapAllocator(), sink)
	handler := newHandler(ModeServer, Config{StartTLS: true})
	if err := ch.Pipeline().AddLast("tls", handler); err != nil {
		t.Fatal(err)
	}
	if inbound != nil {
		if err := ch.Pipeline().AddLast("inbound", inbound); err != nil {
			t.Fatal(err)
		}
	}
	handler.raw = newMemoryConn(handler.bytePool, handler.notifyDrain)
	handler.raw.setNonblocking()
	handler.conn = conn
	handler.startOnce.Do(func() {})
	handler.pendingApplication.markReady()
	handler.started.Store(true)
	handler.handshake = true
	handler.handshakeComplete.Store(true)
	handler.active = true
	handler.activated.Store(true)
	return handler, ch
}

type steadyStateConn struct {
	readData []byte
	written  []byte
	reads    int
}

func (*steadyStateConn) Handshake() error { return nil }

func (*steadyStateConn) ConnectionState() cryptotls.ConnectionState {
	return cryptotls.ConnectionState{}
}

func (c *steadyStateConn) Read(dst []byte) (int, error) {
	c.reads++
	if len(c.readData) == 0 {
		return 0, errWouldBlock
	}
	n := copy(dst, c.readData)
	c.readData = c.readData[n:]
	return n, nil
}

func (c *steadyStateConn) Write(src []byte) (int, error) {
	c.written = append(c.written, src...)
	return len(src), nil
}

func (*steadyStateConn) Close() error { return nil }

var _ Conn = (*steadyStateConn)(nil)

type ciphertextConn struct {
	output *memoryConn
	split  bool
}

func (*ciphertextConn) Handshake() error { return nil }
func (*ciphertextConn) ConnectionState() cryptotls.ConnectionState {
	return cryptotls.ConnectionState{}
}
func (*ciphertextConn) Read([]byte) (int, error) { return 0, errWouldBlock }
func (c *ciphertextConn) Write(src []byte) (int, error) {
	if c.split && len(src) > 1 {
		middle := len(src) / 2
		if _, err := c.output.Write(src[:middle]); err != nil {
			return 0, err
		}
		if _, err := c.output.Write(src[middle:]); err != nil {
			return middle, err
		}
		return len(src), nil
	}
	return c.output.Write(src)
}
func (*ciphertextConn) Close() error { return nil }

type combinedCipherSink struct {
	data           []byte
	writes         int
	flushes        int
	combinedWrites int
}

func (s *combinedCipherSink) Write(message any) error {
	s.writes++
	buf, ok := message.(buffer.ByteBuf)
	if !ok {
		return nil
	}
	s.data = append(s.data, buf.Bytes()...)
	buf.Release()
	return nil
}

func (s *combinedCipherSink) WriteAndFlush(message any) error {
	s.combinedWrites++
	return s.Write(message)
}

func (s *combinedCipherSink) Flush() error {
	s.flushes++
	return nil
}

func (*combinedCipherSink) Close() error { return nil }
