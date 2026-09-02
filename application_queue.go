package tls

import (
	"io"
	"sync"

	"gnalloy.org/gnalloy/buffer"
)

// applicationQueue 只保存握手完成前的应用数据，并在切换稳态时一次性交接所有权。
type applicationQueue struct {
	mu      sync.Mutex
	buffers []buffer.ByteBuf
	ready   bool
	closed  bool
}

func (q *applicationQueue) enqueue(buf buffer.ByteBuf) (bool, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false, io.ErrClosedPipe
	}
	if q.ready {
		return false, nil
	}
	q.buffers = append(q.buffers, buf)
	return true, nil
}

func (q *applicationQueue) markReady() []buffer.ByteBuf {
	q.mu.Lock()
	q.ready = true
	buffers := q.buffers
	q.buffers = nil
	q.mu.Unlock()
	return buffers
}

func (q *applicationQueue) close() []buffer.ByteBuf {
	q.mu.Lock()
	q.closed = true
	buffers := q.buffers
	q.buffers = nil
	q.mu.Unlock()
	return buffers
}

func (q *applicationQueue) len() int {
	q.mu.Lock()
	n := len(q.buffers)
	q.mu.Unlock()
	return n
}
