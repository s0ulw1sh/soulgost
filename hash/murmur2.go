package hash

const (
	murmur2_m_const    = 0x5bd1e995
	murmur2_seed_const = 0
)

func murmur2_mix(h uint32, k uint32) (uint32, uint32) {
	k *= murmur2_m_const
	k ^= k >> 24
	k *= murmur2_m_const
	h *= murmur2_m_const
	h ^= k
	return h, k
}

func MurMur2(data []byte) (h uint32) {
	var k uint32

	h = murmur2_seed_const ^ uint32(len(data))

	for l := len(data); l >= 4; l -= 4 {
		k    = uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
		h, k = murmur2_mix(h, k)
		data = data[4:]
	}

	switch len(data) {
	case 3:
		h ^= uint32(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint32(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint32(data[0])
		h *= murmur2_m_const
	}

	h ^= h >> 13
	h *= murmur2_m_const
	h ^= h >> 15

	return
}
