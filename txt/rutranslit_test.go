package txt

import "testing"

func TestRu2En(t *testing.T) {
	trns := Ru2En("Привет Мир Щек")
	eqls := "Privet Mir Shchek"

	if eqls != trns {
		t.Error("Not equal", eqls, trns)
	}
}