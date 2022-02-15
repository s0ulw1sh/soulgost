package db

import (
	"strconv"
	"encoding/json"
	"database/sql"
	"database/sql/driver"
	"github.com/s0ulw1sh/soulgost/utils"
)

type I32Null struct {
	sql.NullInt32
}

func (n I32Null) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	var out [20]byte
	w := utils.IntToBStr(out[:], int64(n.Int32))
	return out[w:], nil
}

func (n *I32Null) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &n.Int32)
	n.Valid = (err == nil)
	return err
}

func (n *I32Null) Set(val int32) {
	n.Valid  = true
	n.Int32  = val
}

func (n *I32Null) Val() int32 {
	if n.Valid {
		return n.Int32
	} else {
		return 0
	}
}

type I64Null struct {
	sql.NullInt64
}

func (n I64Null) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	var out [20]byte
	w := utils.IntToBStr(out[:], n.Int64)
	return out[w:], nil
}

func (n *I64Null) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &n.Int64)
	n.Valid = (err == nil)
	return err
}

func (n *I64Null) Set(val int64) {
	n.Valid  = true
	n.Int64  = val
}

func (n *I64Null) Val() int64 {
	if n.Valid {
		return n.Int64
	} else {
		return 0
	}
}

type I8Zero struct {
	Int8 int8
}

func (n I8Zero) Value() (driver.Value, error) {
	return int64(n.Int8), nil
}

func (n *I8Zero) Scan(value interface{}) error {
	if value == nil {
		n.Int8 = 0
		return nil
	}

	switch value.(type) {
	case []byte:
		str := string(value.([]byte))
		i64, err := strconv.ParseInt(str, 10, 8)

		if err != nil {
			return err
		}

		n.Int8 = int8(i64)
	case string:
		str := string(value.(string))
		i64, err := strconv.ParseInt(str, 10, 8)

		if err != nil {
			return err
		}

		n.Int8 = int8(i64)
	case int8:   n.Int8 = int8(value.(int8))
	case int16:  n.Int8 = int8(value.(int16))
	case int32:  n.Int8 = int8(value.(int32))
	case int64:  n.Int8 = int8(value.(int64))
	case uint8:  n.Int8 = int8(value.(uint8))
	case uint16: n.Int8 = int8(value.(uint16))
	case uint32: n.Int8 = int8(value.(uint32))
	case uint64: n.Int8 = int8(value.(uint64))
	default:
		return ErrInvalidType
	}

	return nil
}

func (n I8Zero) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Int8)
}

func (n *I8Zero) UnmarshalJSON(b []byte) (err error) {
	if err = json.Unmarshal(b, &n.Int8); err != nil {
		n.Int8 = 0
	}
	return
}

func (n *I8Zero) Set(val int8) {
	n.Int8 = val
}

func (n *I8Zero) Val() int8 {
	return n.Int8
}

type I16Zero struct {
	Var int16
}

func (n I16Zero) Value() (driver.Value, error) {
	return int64(n.Var), nil
}

func (n *I16Zero) Scan(value interface{}) error {
	if value == nil {
		n.Var = 0
		return nil
	}

	switch value.(type) {
	case []byte:
		str := string(value.([]byte))
		i64, err := strconv.ParseInt(str, 10, 16)

		if err != nil {
			return err
		}

		n.Var = int16(i64)
	case string:
		str := string(value.(string))
		i64, err := strconv.ParseInt(str, 10, 16)

		if err != nil {
			return err
		}

		n.Var = int16(i64)
	case int8:   n.Var = int16(value.(int8))
	case int16:  n.Var = int16(value.(int16))
	case int32:  n.Var = int16(value.(int32))
	case int64:  n.Var = int16(value.(int64))
	case uint8:  n.Var = int16(value.(uint8))
	case uint16: n.Var = int16(value.(uint16))
	case uint32: n.Var = int16(value.(uint32))
	case uint64: n.Var = int16(value.(uint64))
	default:
		return ErrInvalidType
	}

	return nil
}

func (n *I16Zero) Set(val int16) {
	n.Var = val
}

func (n *I16Zero) Val() int16 {
	return n.Var
}

func (n I16Zero) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Var)
}

func (n *I16Zero) UnmarshalJSON(b []byte) (err error) {
	if err = json.Unmarshal(b, &n.Var); err != nil {
		n.Var = 0
	}
	return
}

type I32Zero struct {
	Var int32
}

func (n I32Zero) Value() (driver.Value, error) {
	return int64(n.Var), nil
}

func (n *I32Zero) Scan(value interface{}) error {
	if value == nil {
		n.Var = 0
		return nil
	}

	switch value.(type) {
	case []byte:
		str := string(value.([]byte))
		i64, err := strconv.ParseInt(str, 10, 32)

		if err != nil {
			return err
		}

		n.Var = int32(i64)
	case string:
		str := string(value.(string))
		i64, err := strconv.ParseInt(str, 10, 32)

		if err != nil {
			return err
		}

		n.Var = int32(i64)
	case int8:   n.Var = int32(value.(int8))
	case int16:  n.Var = int32(value.(int16))
	case int32:  n.Var = int32(value.(int32))
	case int64:  n.Var = int32(value.(int64))
	case uint8:  n.Var = int32(value.(uint8))
	case uint16: n.Var = int32(value.(uint16))
	case uint32: n.Var = int32(value.(uint32))
	case uint64: n.Var = int32(value.(uint64))
	default:
		return ErrInvalidType
	}

	return nil
}

func (n *I32Zero) Set(val int32) {
	n.Var = val
}

func (n *I32Zero) Val() int32 {
	return n.Var
}

func (n I32Zero) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Var)
}

func (n *I32Zero) UnmarshalJSON(b []byte) (err error) {
	if err = json.Unmarshal(b, &n.Var); err != nil {
		n.Var = 0
	}
	return
}


type I64Zero struct {
	Var int64
}

func (n I64Zero) Value() (driver.Value, error) {
	return int64(n.Var), nil
}

func (n *I64Zero) Scan(value interface{}) error {
	if value == nil {
		n.Var = 0
		return nil
	}

	switch value.(type) {
	case []byte:
		str := string(value.([]byte))
		i64, err := strconv.ParseInt(str, 10, 64)

		if err != nil {
			return err
		}

		n.Var = i64
	case string:
		str := string(value.(string))
		i64, err := strconv.ParseInt(str, 10, 64)

		if err != nil {
			return err
		}

		n.Var = i64
	case int8:   n.Var = int64(value.(int8))
	case int16:  n.Var = int64(value.(int16))
	case int32:  n.Var = int64(value.(int32))
	case int64:  n.Var = int64(value.(int64))
	case uint8:  n.Var = int64(value.(uint8))
	case uint16: n.Var = int64(value.(uint16))
	case uint32: n.Var = int64(value.(uint32))
	case uint64: n.Var = int64(value.(uint64))
	default:
		return ErrInvalidType
	}

	return nil
}

func (n *I64Zero) Set(val int64) {
	n.Var = val
}

func (n *I64Zero) Val() int64 {
	return n.Var
}

func (n I64Zero) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Var)
}

func (n *I64Zero) UnmarshalJSON(b []byte) (err error) {
	if err = json.Unmarshal(b, &n.Var); err != nil {
		n.Var = 0
	}
	return
}