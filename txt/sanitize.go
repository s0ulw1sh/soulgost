package txt

import "bytes"

func Sanitize(text string) string {
	var (
		r   rune
		err error
	)

	var input  = bytes.NewBufferString(text)
	var output = bytes.NewBuffer(nil)

	for {
		if r, _, err = input.ReadRune(); err != nil {
			break
		}

		switch {
		case IsRuRune(r):  fallthrough
		case IsEnRune(r):  fallthrough
		case IsDigitRune(r): fallthrough
		case r == 32:
			output.WriteRune(r)
		}
	}

	return output.String()
}

func SanitizeWithSpace(text string) string {
	var (
		r   rune
		err error
	)

	var input  = bytes.NewBufferString(text)
	var output = bytes.NewBuffer(nil)

	for {
		if r, _, err = input.ReadRune(); err != nil {
			break
		}

		switch {
		case IsRuRune(r):  fallthrough
		case IsEnRune(r):  fallthrough
		case IsDigitRune(r): fallthrough
			output.WriteRune(r)
		}
	}

	return output.String()
}