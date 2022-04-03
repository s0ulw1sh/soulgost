package txt

func IsRuRune(r rune) bool {
	return (r >= 1040 && r <= 1103) || r == 1105 || r == 1025
}

func IsEnRune(r rune) bool {
	return (r >= 65 && r <= 90) || (r >= 97 && r <= 122)
}

func IsDigitRune(r rune) bool {
	return r >= 48 && r <= 57
}
