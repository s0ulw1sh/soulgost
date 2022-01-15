package hash

import "testing"

func TestCRC16(t *testing.T) {
	data := []byte("HelloWorld")
	hash := CRC16Hex(data)
	res  := "7B0A"

	if hash != res {
		t.Error("Not equal hashes", hash, res)
	}
}