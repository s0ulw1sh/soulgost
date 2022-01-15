package utils

import "testing"

func TestIntToBStr(t *testing.T) {
	var (
		w    int
		ns   string
		bstr [9]byte
		nums map[string]int64 = map[string]int64{
			"0":        0,
			"123":      123,
			"-222":     -222,
			"-1":       -1,
			"-9999":    -9999,
			"-701":     -701,
			"3912341":  3912341,
		}
	)

	for s, n := range nums {
		w = IntToBStr(bstr[:], n)

		if ns = string(bstr[w:]); ns != s {
			t.Error("Not equal strings", ns, s)
		}
	}
}