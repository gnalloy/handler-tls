package tls

import (
	"io"
	"sync"
	"sync/atomic"

	"gnalloy.org/gnalloy/buffer"
)

const (
	applicationQueuePending uint32 = iota
	applicationQueueReady
	applicationQueueClosed
)

// applicationQueue 只保存握手完成前的应用数据，并在切换稳态时一次性交接所有权。
type applicationQueue struct {
	mu      sync.Mutex
	buffers []buffer.ByteBuf
	state   atomic.Uint32
}

func (q *applicationQueue) enqueue(buf buffer.ByteBuf) (bool, error) {
	switch q.state.Load() {
	case applicationQueueReady:
		return false, nil
	case applicationQueueClosed:
		return false, io.ErrClosedPipe
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	switch q.state.Load() {
	case applicationQueueReady:
		return false, nil
	case applicationQueueClosed:
		return false, io.ErrClosedPipe
	}
	q.buffers = append(q.buffers, buf)
	return true, nil
}

func (q *applicationQueue) markReady() []buffer.ByteBuf {
	q.mu.Lock()
	if q.state.Load() == applicationQueueClosed {
		q.mu.Unlock()
		return nil
	}
	buffers := q.buffers
	q.buffers = nil
	q.state.Store(applicationQueueReady)
	q.mu.Unlock()
	return buffers
}

func (q *applicationQueue) close() []buffer.ByteBuf {
	q.mu.Lock()
	buffers := q.buffers
	q.buffers = nil
	q.state.Store(applicationQueueClosed)
	q.mu.Unlock()
	return buffers
}

func (q *applicationQueue) len() int {
	q.mu.Lock()
	n := len(q.buffers)
	q.mu.Unlock()
	return n
}
