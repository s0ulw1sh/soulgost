package types

import (
	"database/sql"
	"encoding/json"
)

type F64Null struct {
	sql.NullFloat64
}

func (f *F64Null) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(f.Float64)
}

func (f *F64Null) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &f.Float64)
	f.Valid = (err == nil)
	return err
}

func (n *F64Null) Set(val float64) {
	n.Valid    = true
	n.Float64  = val
}

func (n *F64Null) Val() float64 {
	if n.Valid {
		return n.Float64
	} else {
		return 0
	}
}