package tls

import (
	"errors"
	"io"
	"testing"
)

func TestApplicationQueueRejectsEnqueueAfterClose(t *testing.T) {
	queue := &applicationQueue{}
	queue.markReady()
	queue.close()

	queued, err := queue.enqueue(nil)
	if queued || !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("enqueue after close=(%t, %v), want (false, io.ErrClosedPipe)", queued, err)
	}
}

func BenchmarkApplicationQueueReadyEnqueue(b *testing.B) {
	queue := &applicationQueue{}
	queue.markReady()
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		queued, err := queue.enqueue(nil)
		if queued || err != nil {
			b.Fatalf("enqueue=(%t, %v), want (false, nil)", queued, err)
		}
	}
}
