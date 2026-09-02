package tls

import "testing"

func TestOwnedBufferFromChunkTransfersPoolOwnership(t *testing.T) {
	pool := &trackingBytePool{}
	chunk := newByteChunk(pool.Acquire(128), pool)
	out := ownedBufferFromChunk(&chunk)
	if out == nil {
		t.Fatal("owned buffer is nil")
	}
	if chunk.data != nil || chunk.pool != nil {
		t.Fatal("chunk retained ownership after transfer")
	}
	if !out.Release() {
		t.Fatal("owned buffer release did not drop final reference")
	}
	if pool.released != 1 {
		t.Fatalf("releases=%d, want 1", pool.released)
	}
}

func TestOwnedBufferFromChunkRoundTripDoesNotAllocate(t *testing.T) {
	pool := NewPooledBytePool(BytePoolConfig{})
	chunk := newByteChunk(pool.Acquire(128), pool)
	ownedBufferFromChunk(&chunk).Release()

	allocs := testing.AllocsPerRun(1000, func() {
		chunk := newByteChunk(pool.Acquire(128), pool)
		ownedBufferFromChunk(&chunk).Release()
	})
	if allocs != 0 {
		t.Fatalf("allocs=%f, want 0", allocs)
	}
}

func BenchmarkOwnedBufferFromChunkRoundTrip(b *testing.B) {
	pool := NewPooledBytePool(BytePoolConfig{})
	chunk := newByteChunk(pool.Acquire(128), pool)
	ownedBufferFromChunk(&chunk).Release()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		chunk := newByteChunk(pool.Acquire(128), pool)
		ownedBufferFromChunk(&chunk).Release()
	}
}
