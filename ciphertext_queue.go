package tls

import "sync"

// ciphertextQueue 在握手协程与 Channel owner EventLoop 之间转移密文所有权。
type ciphertextQueue struct {
	mu     sync.Mutex
	chunks []byteChunk
	head   int
	closed bool
}

func (q *ciphertextQueue) push(chunk byteChunk) bool {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return false
	}
	q.chunks = append(q.chunks, chunk)
	q.mu.Unlock()
	return true
}

func (q *ciphertextQueue) pop() (chunk byteChunk, remaining bool, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.head == len(q.chunks) {
		return byteChunk{}, false, false
	}
	chunk = q.chunks[q.head]
	q.chunks[q.head] = byteChunk{}
	q.head++
	remaining = q.head != len(q.chunks)
	if q.head == len(q.chunks) {
		q.chunks = q.chunks[:0]
		q.head = 0
	}
	return chunk, remaining, true
}

func (q *ciphertextQueue) len() int {
	q.mu.Lock()
	n := len(q.chunks) - q.head
	q.mu.Unlock()
	return n
}

func (q *ciphertextQueue) closeAndRelease() {
	q.mu.Lock()
	chunks := q.chunks
	head := q.head
	q.closed = true
	q.chunks = nil
	q.head = 0
	q.mu.Unlock()
	for i := head; i < len(chunks); i++ {
		chunks[i].releaseOwned()
	}
}
