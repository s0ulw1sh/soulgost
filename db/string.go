package db

import (
	"database/sql"
	"encoding/json"
)

type StrNull struct {
	sql.NullString
}

func StrNullNew(str string) StrNull {
	var s StrNull
	s.Set(str)
	return s
}

func (s *StrNull) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (s *StrNull) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &s.String)
	s.Valid = (err == nil)
	return err
}

func (s *StrNull) Set(val string) {
	s.Valid  = val != ""
	s.String = val
}

func (s *StrNull) Val() string {
	if s.Valid {
		return s.String
	} else {
		return ""
	}
}

type StrEmpty struct {
	sql.NullString
}

func (s *StrEmpty) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("\"-\""), nil
	}
	return json.Marshal(s.String)
}

func (s *StrEmpty) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &s.String)
	s.Valid = (err == nil)
	return err
}

func (s *StrEmpty) Set(val string) {
	s.Valid  = val != ""
	s.String = val
}

func (s *StrEmpty) Val() string {
	if s.Valid {
		return s.String
	} else {
		return ""
	}
}