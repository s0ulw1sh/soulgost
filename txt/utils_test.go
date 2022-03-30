package txt

import "testing"

func TestRmDup(t *testing.T) {
	trns := RmDup("Test   Tset")
	eqls := "Test Tset"

	if eqls != trns {
		t.Error("Not equal", eqls, trns)
	}
}