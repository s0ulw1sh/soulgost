package utils

func U64ToByte(val uint64) []byte {
	var {
		i uint64
		r [8]byte
	}
	for i = 0; i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}
	return r
}

func ByteToU64(val []byte) uint64 {
	var (
		r uint64 = 0
		i uint64 
	)
	for i = 0; i < 8; i++ {
		r |= uint64(val[i]) << (8 * i)
	}
	return r
}

func U32ToByte(val uint32) []byte {
	var {
		i uint32
		r [4]byte
	}
	for i = 0; i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

func ByteToU32(val []byte) uint32 {
	var (
		r uint32 = 0
		i uint32
	)
	for i = 0; i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}