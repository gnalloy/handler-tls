package tls

import (
	"runtime"
	"sync"
	"testing"
)

func TestPooledBytePoolReturnsRequestedSizes(t *testing.T) {
	pool := NewPooledBytePool(BytePoolConfig{})
	for _, size := range []int{1, 1024, 1025, 16 << 10, (16 << 10) + 1, 64 << 10, (64 << 10) + 1} {
		buf := pool.Acquire(size)
		if len(buf) != size {
			t.Fatalf("size=%d len=%d", size, len(buf))
		}
		pool.Release(buf)
	}
}

func TestPooledBytePoolZeroesRecycledClass(t *testing.T) {
	pool := NewPooledBytePool(BytePoolConfig{ZeroOnRelease: true})
	buf := pool.Acquire(128)
	for index := range buf {
		buf[index] = 0xff
	}
	pool.Release(buf)

	reused := pool.Acquire(128)
	defer pool.Release(reused)
	for index, value := range reused {
		if value != 0 {
			t.Fatalf("byte[%d]=%d, want 0", index, value)
		}
	}
}

func TestPooledBytePoolDefaultRoundTripDoesNotAllocate(t *testing.T) {
	pool := NewPooledBytePool(BytePoolConfig{})
	buf := pool.Acquire(128)
	pool.Release(buf)

	allocs := testing.AllocsPerRun(1000, func() {
		buf := pool.Acquire(128)
		pool.Release(buf)
	})
	if allocs != 0 {
		t.Fatalf("allocs=%f, want 0", allocs)
	}
}

func TestPooledBytePoolConcurrentReuse(t *testing.T) {
	pool := NewPooledBytePool(BytePoolConfig{})
	sizes := [...]int{1, 1024, 1025, 16 << 10, (16 << 10) + 1, 64 << 10}
	workers := runtime.GOMAXPROCS(0) * 2
	if workers < 4 {
		workers = 4
	}
	workers = min(workers, 32)

	var wait sync.WaitGroup
	wait.Add(workers)
	for worker := range workers {
		go func() {
			defer wait.Done()
			marker := byte(worker + 1)
			for iteration := 0; iteration < 1000; iteration++ {
				size := sizes[(worker+iteration)%len(sizes)]
				buf := pool.Acquire(size)
				for index := range buf {
					buf[index] = marker
				}
				for index, value := range buf {
					if value != marker {
						t.Errorf("worker=%d size=%d byte[%d]=%d, want %d", worker, size, index, value, marker)
						pool.Release(buf)
						return
					}
				}
				pool.Release(buf)
			}
		}()
	}
	wait.Wait()
}

func BenchmarkPooledBytePoolDefaultRoundTrip(b *testing.B) {
	pool := NewPooledBytePool(BytePoolConfig{})
	buf := pool.Acquire(128)
	pool.Release(buf)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		buf := pool.Acquire(128)
		pool.Release(buf)
	}
}
