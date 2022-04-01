package txt

import "bytes"

var ruEnDictBase = map[string]string{
	"а": "a", "А": "A", "Б": "B", "б": "b", "В": "V", "в": "v", "Г": "G", "г": "g",
	"Д": "D", "д": "d", "З": "Z", "з": "z", "И": "I", "и": "i", "К": "K", "к": "k",
	"Л": "L", "л": "l", "М": "M", "м": "m", "Н": "N", "н": "n", "О": "O", "о": "o",
	"П": "P", "п": "p", "Р": "R", "р": "r", "С": "S", "с": "s", "Т": "T", "т": "t",
	"У": "U", "у": "u", "Ф": "F", "ф": "f",
}

var ruEnDictExt = map[string]string{
	"Е": "E", "е": "e", "Ё": "E", "ё": "e", "Ж": "Zh", "ж": "zh", "Й": "I", "й": "i",
	"Х": "Kh", "х": "kh", "Ц": "Ts", "ц": "ts", "Ч": "Ch", "ч": "ch", "Ш": "Sh",
	"ш": "sh", "Щ": "Shch", "щ": "shch", "Ъ": "Ie", "ъ": "ie", "Ы": "Y", "ы": "y",
	"Ь": "", "ь": "", "Э": "E", "э": "e", "Ю": "Iu", "ю": "iu", "Я": "Ia", "я": "ia",
}

func IsRuChar(r rune) bool {
	return (r >= 1040 && r <= 1103) || r == 1105 || r == 1025
}

func Ru2En(text string) string {
	var (
		r   rune
		rr  string
		ok  bool
		err error
	)

	var input  = bytes.NewBufferString(text)
	var output = bytes.NewBuffer(nil)

	for {
		if r, _, err = input.ReadRune(); err != nil {
			break
		} else if !IsRuChar(r) {
			output.WriteRune(r)
			continue
		}

		if rr, ok = ruEnDictBase[string(r)]; ok {
			output.WriteString(rr)
		} else if rr, ok = ruEnDictExt[string(r)]; ok {
			output.WriteString(rr)
		}
	}

	return output.String()
}