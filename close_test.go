package tls

import (
	"testing"

	"gnalloy.org/gnalloy/buffer"
)

func TestCloseReleasesQueuedApplicationBuffers(t *testing.T) {
	released := 0
	handler := newHandler(ModeServer, Config{})
	buf := buffer.NewOwnedBuffer([]byte("queued"), func([]byte) {
		released++
	})
	queued, err := handler.pendingApplication.enqueue(buf)
	if err != nil || !queued {
		t.Fatalf("enqueue queued=%v err=%v", queued, err)
	}

	handler.close()

	if released != 1 {
		t.Fatalf("released=%d, want 1", released)
	}
	if pending := handler.pendingApplication.len(); pending != 0 {
		t.Fatalf("pending=%d, want 0", pending)
	}
}
