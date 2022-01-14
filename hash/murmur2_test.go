package hash

import "testing"

func TestMurMur2(t *testing.T) {
	data := []byte("HelloWorld")
	hash := MurMur2(data)
	eq   := uint32(147875502)

	if hash != eq {
		t.Error("Not equal hashes", hash, eq)
	}
}