package tls

import "testing"

func TestCiphertextQueuePopReportsRemaining(t *testing.T) {
	queue := &ciphertextQueue{}
	pool := &trackingBytePool{}
	if !queue.push(newByteChunk([]byte("first"), pool)) {
		t.Fatal("first push failed")
	}
	if !queue.push(newByteChunk([]byte("second"), pool)) {
		t.Fatal("second push failed")
	}

	first, remaining, ok := queue.pop()
	if !ok || !remaining || string(first.data) != "first" {
		t.Fatalf("first pop=(%q, %t, %t), want (first, true, true)", first.data, remaining, ok)
	}
	first.releaseOwned()

	second, remaining, ok := queue.pop()
	if !ok || remaining || string(second.data) != "second" {
		t.Fatalf("second pop=(%q, %t, %t), want (second, false, true)", second.data, remaining, ok)
	}
	second.releaseOwned()

	_, remaining, ok = queue.pop()
	if ok || remaining {
		t.Fatalf("empty pop=(%t, %t), want (false, false)", remaining, ok)
	}
	if pool.released != 2 {
		t.Fatalf("released=%d, want 2", pool.released)
	}
}
