package hash

type MurMur2Hash struct {
	m_const uint32
	hash    uint32
	lcount  int
	last    [3]byte
}

func (self *MurMur2Hash) mix(h uint32, k uint32) (uint32, uint32) {
	k *= self.m_const
	k ^= k >> 24
	k *= self.m_const
	h *= self.m_const
	h ^= k

	return h, k
}

func (self *MurMur2Hash) Init(seed uint32) {
	self.m_const = 0x5bd1e995
	self.hash    = seed
}

func (self *MurMur2Hash) Update(data []byte) {
	var k uint32
	var l int = len(data)

	if l == 0 { return }

	if self.lcount > 0 && l - self.lcount > 0 {
		k = 0

		switch self.lcount {
		case 3:
			k = uint32(self.last[0]) | uint32(self.last[1])<<8 | uint32(self.last[2])<<16 | uint32(data[0])<<24
		case 2:
			k = uint32(self.last[0]) | uint32(self.last[1])<<8 | uint32(data[0])<<16 | uint32(data[1])<<24
		case 1:
			k = uint32(self.last[0]) | uint32(data[0])<<8 | uint32(data[1])<<16 | uint32(data[2])<<24
		}

		l = l - (4-self.lcount)
		data = data[4-self.lcount:]
		self.hash, k = self.mix(self.hash, k)
		self.lcount  = 0
	}

	for ; l >= 4; l -= 4 {
		k            = uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
		self.hash, k = self.mix(self.hash, k)
		data         = data[4:]
	}

	self.lcount = len(data)

	switch self.lcount {
	case 3:
		self.last[2] = data[2]
		fallthrough
	case 2:
		self.last[1] = data[1]
		fallthrough
	case 1:
		self.last[0] = data[0]
	}
}

func (self *MurMur2Hash) Finish() uint32 {

	switch self.lcount {
	case 3:
		self.hash ^= uint32(self.last[2]) << 16
		fallthrough
	case 2:
		self.hash ^= uint32(self.last[1]) << 8
		fallthrough
	case 1:
		self.hash ^= uint32(self.last[0])
		self.hash *= self.m_const
	}

	self.hash ^= self.hash >> 13
	self.hash *= self.m_const
	self.hash ^= self.hash >> 15

	return self.hash
}

func MurMur2(data []byte) uint32 {
	var h MurMur2Hash

	h.Init(uint32(0 ^ len(data)))
	h.Update(data)

	return h.Finish()
}
