package tls

import "unsafe"

const bytePoolClassCount = 7

var bytePoolClassSizes = [bytePoolClassCount]int{
	1 << 10,
	2 << 10,
	4 << 10,
	8 << 10,
	16 << 10,
	32 << 10,
	64 << 10,
}

type (
	byteBlock1K  [1 << 10]byte
	byteBlock2K  [2 << 10]byte
	byteBlock4K  [4 << 10]byte
	byteBlock8K  [8 << 10]byte
	byteBlock16K [16 << 10]byte
	byteBlock32K [32 << 10]byte
	byteBlock64K [64 << 10]byte
)

func bytePoolClass(size int, maxSize int) int {
	for index, classSize := range bytePoolClassSizes {
		if classSize >= size && classSize <= maxSize {
			return index
		}
	}
	return -1
}

func bytePoolExactClass(capacity int) int {
	for index, classSize := range bytePoolClassSizes {
		if classSize == capacity {
			return index
		}
	}
	return -1
}

func (p *PooledBytePool) acquireClass(index int, size int) []byte {
	value := p.pools[index].Get()
	switch index {
	case 0:
		if value == nil {
			value = new(byteBlock1K)
		}
		return value.(*byteBlock1K)[:size]
	case 1:
		if value == nil {
			value = new(byteBlock2K)
		}
		return value.(*byteBlock2K)[:size]
	case 2:
		if value == nil {
			value = new(byteBlock4K)
		}
		return value.(*byteBlock4K)[:size]
	case 3:
		if value == nil {
			value = new(byteBlock8K)
		}
		return value.(*byteBlock8K)[:size]
	case 4:
		if value == nil {
			value = new(byteBlock16K)
		}
		return value.(*byteBlock16K)[:size]
	case 5:
		if value == nil {
			value = new(byteBlock32K)
		}
		return value.(*byteBlock32K)[:size]
	case 6:
		if value == nil {
			value = new(byteBlock64K)
		}
		return value.(*byteBlock64K)[:size]
	default:
		return nil
	}
}

func (p *PooledBytePool) releaseClass(index int, buf []byte) {
	// buf 只能来自 acquireClass，首元素地址就是固定容量块的起始地址。
	pointer := unsafe.Pointer(&buf[:1][0])
	switch index {
	case 0:
		p.pools[index].Put((*byteBlock1K)(pointer))
	case 1:
		p.pools[index].Put((*byteBlock2K)(pointer))
	case 2:
		p.pools[index].Put((*byteBlock4K)(pointer))
	case 3:
		p.pools[index].Put((*byteBlock8K)(pointer))
	case 4:
		p.pools[index].Put((*byteBlock16K)(pointer))
	case 5:
		p.pools[index].Put((*byteBlock32K)(pointer))
	case 6:
		p.pools[index].Put((*byteBlock64K)(pointer))
	}
}
