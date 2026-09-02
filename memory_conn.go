package tls

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var errWouldBlock error = errWouldBlockError{}

type errWouldBlockError struct{}

func (errWouldBlockError) Error() string   { return "gnalloy/handler/tls: operation would block" }
func (errWouldBlockError) Timeout() bool   { return false }
func (errWouldBlockError) Temporary() bool { return true }

func (errWouldBlockError) Is(target error) bool {
	_, ok := target.(errWouldBlockError)
	return ok
}

type memoryConn struct {
	inMu        sync.Mutex
	inCond      *sync.Cond
	in          []byteChunk
	inHead      int
	closed      bool
	once        sync.Once
	pool        BytePool
	notify      func()
	pending     byteChunk
	out         ciphertextQueue
	nonblocking atomic.Bool
}

func newMemoryConn(pool BytePool, notify func()) *memoryConn {
	c := &memoryConn{
		pool:   normalizeBytePool(pool),
		notify: notify,
	}
	c.inCond = sync.NewCond(&c.inMu)
	return c
}

func (c *memoryConn) feed(src []byte) error {
	if len(src) == 0 {
		return nil
	}
	return c.feedOwned(copyBytes(src, c.pool))
}

func (c *memoryConn) feedOwned(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	chunk := newByteChunk(data, c.pool)
	c.inMu.Lock()
	if c.closed {
		c.inMu.Unlock()
		chunk.releaseOwned()
		return io.ErrClosedPipe
	}
	c.in = append(c.in, chunk)
	c.inCond.Signal()
	c.inMu.Unlock()
	return nil
}

func (c *memoryConn) Read(dst []byte) (int, error) {
	c.inMu.Lock()
	defer c.inMu.Unlock()
	for len(c.pending.data) == 0 {
		if c.inHead < len(c.in) {
			c.pending = c.in[c.inHead]
			c.in[c.inHead] = byteChunk{}
			c.inHead++
			if c.inHead == len(c.in) {
				c.in = c.in[:0]
				c.inHead = 0
			}
			continue
		}
		if c.closed {
			return 0, io.EOF
		}
		if c.nonblocking.Load() {
			return 0, errWouldBlock
		}
		c.inCond.Wait()
	}
	n := copy(dst, c.pending.data)
	c.pending.data = c.pending.data[n:]
	if len(c.pending.data) == 0 {
		c.pending.releaseOwned()
	}
	return n, nil
}

func (c *memoryConn) Write(src []byte) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}
	chunk := newByteChunk(copyBytes(src, c.pool), c.pool)
	if !c.out.push(chunk) {
		chunk.releaseOwned()
		return 0, io.ErrClosedPipe
	}
	c.notifyDrain()
	return len(src), nil
}

func (c *memoryConn) Close() error {
	c.once.Do(func() {
		c.inMu.Lock()
		c.closed = true
		c.pending.releaseOwned()
		queued := c.in
		head := c.inHead
		c.in = nil
		c.inHead = 0
		c.inCond.Broadcast()
		c.inMu.Unlock()
		for i := head; i < len(queued); i++ {
			queued[i].releaseOwned()
		}
	})
	return nil
}

func (c *memoryConn) setNonblocking() {
	c.nonblocking.Store(true)
}

func (c *memoryConn) popOutput() (byteChunk, bool) {
	return c.out.pop()
}

func (c *memoryConn) hasOutput() bool {
	return c.out.len() != 0
}

func (c *memoryConn) releaseOutput() {
	c.out.closeAndRelease()
}

func (c *memoryConn) LocalAddr() net.Addr              { return memoryAddr("local") }
func (c *memoryConn) RemoteAddr() net.Addr             { return memoryAddr("remote") }
func (c *memoryConn) SetDeadline(time.Time) error      { return nil }
func (c *memoryConn) SetReadDeadline(time.Time) error  { return nil }
func (c *memoryConn) SetWriteDeadline(time.Time) error { return nil }

func (c *memoryConn) notifyDrain() {
	if c.notify != nil {
		c.notify()
	}
}

type memoryAddr string

func (a memoryAddr) Network() string { return "memory" }
func (a memoryAddr) String() string  { return string(a) }

var _ net.Error = errWouldBlockError{}
