package utils

func IntToBStr(buf []byte, v int64) int {
	w := len(buf)
	m := v >> 63

	if v == 0 {
		w--
		buf[w] = '0'
		return w
	}

	for v != 0 {
		w--
		buf[w] = byte((v%10 ^ m) - m) + '0'
		v /= 10
	}

	if m < 0 {
		w--
		buf[w] = '-'
	}
	
	return w
}

func UintToBStr(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

func UintToBStrLeadZero(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
		w--
		buf[w] = '0'
	} else {
		ism := v < 10
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
		if ism {
			w--
			buf[w] = '0'
		}
	}
	return w
}

