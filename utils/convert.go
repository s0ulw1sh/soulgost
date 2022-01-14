package utils

/*func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}*/

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

