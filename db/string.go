package db

import (
	"database/sql"
	"encoding/json"
	"database/sql/driver"
)

type StrNull struct {
	ns sql.NullString
}

func StrNullNew(str string) StrNull {
	var s StrNull
	s.Set(str)
	return s
}

func (s StrNull) String() string {
	if !s.ns.Valid {
		return ""
	}
	return s.ns.String
}

func (s *StrNull) Scan(value interface{}) error {
	return s.ns.Scan(value)
}

func (s StrNull) Value() (driver.Value, error) {
	return s.ns.Value()
}

func (s StrNull) MarshalJSON() ([]byte, error) {
	if !s.ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.ns.String)
}

func (s *StrNull) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &s.ns.String)
	s.ns.Valid = (err == nil)
	return err
}

func (s *StrNull) Set(val string) {
	s.ns.Valid  = val != ""
	s.ns.String = val
}

func (s *StrNull) Val() string {
	if s.ns.Valid {
		return s.ns.String
	} else {
		return ""
	}
}

type StrEmpty struct {
	ns sql.NullString
}

func (s *StrEmpty) Scan(value interface{}) error {
	return s.ns.Scan(value)
}

func (s StrEmpty) Value() (driver.Value, error) {
	return s.ns.Value()
}

func (s StrEmpty) String() string {
	if !s.ns.Valid {
		return "-"
	}
	return s.ns.String
}

func (s StrEmpty) MarshalJSON() ([]byte, error) {
	if !s.ns.Valid {
		return []byte("\"-\""), nil
	}
	return json.Marshal(s.ns.String)
}

func (s *StrEmpty) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &s.ns.String)
	s.ns.Valid = (err == nil)
	return err
}

func (s *StrEmpty) Set(val string) {
	s.ns.Valid  = val != ""
	s.ns.String = val
}

func (s *StrEmpty) Val() string {
	if s.ns.Valid {
		return s.ns.String
	} else {
		return ""
	}
}