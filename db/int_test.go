package db

import "testing"
import "encoding/json"


func TestInt64(t *testing.T) {
	var i I64Zero

	i.Set(22)

	b, err := json.Marshal(i)

	if err != nil || string(b) != "22" {
		t.Error("Invalid Marshaling I64Zero", string(b), err)
	}
}