package tls

import (
	"errors"
	"net"
	"testing"

	"gnalloy.org/gnalloy/buffer"
)

func TestMemoryConnReleasesInputAfterRead(t *testing.T) {
	pool := &trackingBytePool{}
	conn := newMemoryConn(pool, nil)
	data := copyBytes([]byte("abcd"), pool)
	if err := conn.feedOwned(data); err != nil {
		t.Fatal(err)
	}

	var dst [2]byte
	if n, err := conn.Read(dst[:]); err != nil || n != 2 || string(dst[:]) != "ab" {
		t.Fatalf("first read n=%d err=%v dst=%q", n, err, string(dst[:]))
	}
	if pool.released != 0 {
		t.Fatalf("released=%d before chunk is fully consumed", pool.released)
	}
	if n, err := conn.Read(dst[:]); err != nil || n != 2 || string(dst[:]) != "cd" {
		t.Fatalf("second read n=%d err=%v dst=%q", n, err, string(dst[:]))
	}
	if pool.released != 1 {
		t.Fatalf("released=%d, want input chunk released", pool.released)
	}
}

func TestMemoryConnTakesByteBufOwnershipWithoutCopy(t *testing.T) {
	conn := newMemoryConn(nil, nil)
	buf := buffer.NewHeapBuffer(4)
	_, _ = buf.WriteBytes([]byte("data"))
	if err := conn.feedBuffer(buf); err != nil {
		t.Fatal(err)
	}
	if buf.RefCnt() != 1 {
		t.Fatalf("ref=%d before read, want 1", buf.RefCnt())
	}

	var dst [4]byte
	if n, err := conn.Read(dst[:]); err != nil || n != 4 || string(dst[:]) != "data" {
		t.Fatalf("read n=%d err=%v data=%q", n, err, string(dst[:]))
	}
	if buf.RefCnt() != 0 {
		t.Fatalf("ref=%d after read, want released", buf.RefCnt())
	}
}

func TestMemoryConnCloseReleasesQueuedInput(t *testing.T) {
	pool := &trackingBytePool{}
	conn := newMemoryConn(pool, nil)
	data := copyBytes([]byte("cipher"), pool)
	if err := conn.feedOwned(data); err != nil {
		t.Fatal(err)
	}
	if pool.released != 0 {
		t.Fatalf("released=%d before close", pool.released)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	if pool.released != 1 {
		t.Fatalf("released=%d, want queued input released", pool.released)
	}
}

func TestMemoryConnNotifiesWhenCiphertextIsReady(t *testing.T) {
	notified := make(chan struct{}, 1)
	conn := newMemoryConn(nil, func() {
		select {
		case notified <- struct{}{}:
		default:
		}
	})
	if _, err := conn.Write([]byte("cipher")); err != nil {
		t.Fatal(err)
	}
	select {
	case <-notified:
	default:
		t.Fatal("ciphertext write did not notify drain")
	}
}

func TestMemoryConnReturnsTemporaryErrorInNonblockingMode(t *testing.T) {
	conn := newMemoryConn(nil, nil)
	conn.setNonblocking()

	var dst [1]byte
	_, err := conn.Read(dst[:])
	if !errors.Is(err, errWouldBlock) {
		t.Fatalf("err=%v, want errWouldBlock", err)
	}
	var netErr net.Error
	if !errors.As(err, &netErr) || !netErr.Temporary() || netErr.Timeout() {
		t.Fatalf("err=%v, want temporary non-timeout network error", err)
	}
}

func TestEvaluateNativeProviderRequiresTLS13ALPNAndQUIC(t *testing.T) {
	evaluation := EvaluateNativeProvider(nativeProviderStub{
		Provider:             "stub",
		TLS13:                true,
		ALPN:                 true,
		QUICPacketProtection: true,
	})
	if !evaluation.Supported {
		t.Fatalf("evaluation=%+v, want supported", evaluation)
	}

	evaluation = EvaluateNativeProvider(UnsupportedNativeProvider{})
	if evaluation.Supported || len(evaluation.Reasons) == 0 {
		t.Fatalf("evaluation=%+v, want unsupported reasons", evaluation)
	}
}

type nativeProviderStub NativeCapabilities

func (p nativeProviderStub) Capabilities() NativeCapabilities {
	return NativeCapabilities(p)
}

type trackingBytePool struct {
	acquired int
	released int
}

func (p *trackingBytePool) Acquire(size int) []byte {
	p.acquired++
	return make([]byte, size)
}

func (p *trackingBytePool) Release([]byte) {
	p.released++
}
