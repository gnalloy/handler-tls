package tls

func copyBytes(src []byte, pool BytePool) []byte {
	if len(src) == 0 {
		return nil
	}
	out := acquireBytes(pool, len(src))
	copy(out, src)
	return out
}
