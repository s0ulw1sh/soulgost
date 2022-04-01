package utils

import "testing"

func TestToByte(t *testing.T) {
	a := U32ToByte(127)
	b := U64ToByte(123)
	y := [4]byte{127, 0, 0, 0}
	z := [8]byte{123, 0, 0, 0, 0, 0, 0, 0}

	if  a != y {
		t.Error("Not equal", a, y)
	}

	if  b != z {
		t.Error("Not equal", b, z)
	}

	if  ByteToU64(b) != 123 {
		t.Error("Not equal")
	}

	if  ByteToU32(a) != 127 {
		t.Error("Not equal")
	}
}