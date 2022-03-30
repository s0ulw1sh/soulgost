package txt

func RmDup(word string) string {
	previousChar := []rune(word)[0]
	result := string(previousChar)
	for _, rune := range word[1:] {
		if rune != previousChar || rune == 'C' {
			result = result + string(rune)
		}
		previousChar = rune
	}
	return result
}

func SubStr(s string, start, count int) string {
	l := len(s)
	if start < 0 {
		start = 0
		count = count + start
	}
	if start+count > l {
		count = l - start
	}
	return s[start : start+count]
}