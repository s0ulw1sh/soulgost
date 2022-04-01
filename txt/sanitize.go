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

		if r >= 33 && r <= 47 {
			continue
		} else if r >= 58 && r <= 64 {
			continue
		} else if r >= 91 && r <= 96 {
			continue
		} else if r >= 123 && r <= 126 {
			continue
		} else if r >= 8 && r <= 10 {
			continue
		}

		output.WriteRune(r)
	}

	return output.String()
}