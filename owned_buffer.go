package tls

import (
	"gnalloy.org/gnalloy/buffer"
)

func ownedBufferFromChunk(chunk *byteChunk) buffer.ByteBuf {
	if chunk == nil || len(chunk.data) == 0 {
		return nil
	}
	data := chunk.data
	pool := chunk.pool
	*chunk = byteChunk{}
	return buffer.NewOwnedBufferWithReleaser(data, pool)
}
