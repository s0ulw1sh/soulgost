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

func TestMurMur2Part(t *testing.T) {
	data1 := []byte("Hello")
	data2 := []byte("World")

	var h MurMur2Hash

	h.Init(10)
	h.Update(data1)
	h.Update(data2)

	hash := h.Finish()
	eq   := uint32(147875502)

	if hash != eq {
		t.Error("Not equal hashes", hash, eq)
	}
}