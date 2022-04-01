package utils

const u64ToByteMask = 0x00000000000000FF
const u32ToByteMask = 0x000000FF

func U64ToByte(val uint64) (r [8]byte) {
	r[0] = byte(val & u64ToByteMask)
	val  = val >> 8
	r[1] = byte(val & u64ToByteMask)
	val  = val >> 8
	r[2] = byte(val & u64ToByteMask)
	val  = val >> 8
	r[3] = byte(val & u64ToByteMask)
	val  = val >> 8
	r[4] = byte(val & u64ToByteMask)
	val  = val >> 8
	r[5] = byte(val & u64ToByteMask)
	val  = val >> 8
	r[6] = byte(val & u64ToByteMask)
	val  = val >> 8
	r[7] = byte(val & u64ToByteMask)

	return
}

func ByteToU64(val [8]byte) (r uint64) {
	r |= uint64(val[0])
	r |= uint64(val[1]) << 8
	r |= uint64(val[2]) << 16
	r |= uint64(val[3]) << 24
	r |= uint64(val[4]) << 32
	r |= uint64(val[5]) << 40
	r |= uint64(val[6]) << 48
	r |= uint64(val[7]) << 56

	return
}

func U32ToByte(val uint32) (r [4]byte) {
	r[0] = byte(val & u32ToByteMask)
	val  = val >> 8
	r[1] = byte(val & u32ToByteMask)
	val  = val >> 8
	r[2] = byte(val & u32ToByteMask)
	val  = val >> 8
	r[3] = byte(val & u32ToByteMask)

	return
}

func ByteToU32(val [4]byte) (r uint32) {
	r |= uint32(val[0])
	r |= uint32(val[1]) << 8
	r |= uint32(val[2]) << 16
	r |= uint32(val[3]) << 24
	return r
}