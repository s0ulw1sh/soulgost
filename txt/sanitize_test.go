package txt

import "testing"

func TestSanitize(t *testing.T) {
	snts := Sanitize("`Hello` - Wo+r-ld! «Привет»‎ % Ы@Р	X_©")
	eqls := "Hello  World Привет  ЫРX"

	if eqls != snts {
		t.Error("Not equal", eqls, "|", snts)
	}
}