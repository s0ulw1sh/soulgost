package txt

import "testing"

func TestMetaphone(t *testing.T) {
	trns := Metaphone("TEST")
	eqls := "TST"

	if eqls != trns {
		t.Error("Not equal", eqls, trns)
	}
}