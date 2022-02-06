package db

import (
    "time"
    "strings"
    "strconv"
    "encoding/json"
    "database/sql/driver"
    "github.com/s0ulw1sh/soulgost/utils"
)

type TimeMin struct {
    Minutes int
}

func (d TimeMin) Value() (driver.Value, error) {
    return d.String(), nil
}

func (d *TimeMin) Scan(value interface{}) error {
    var str string

    switch value.(type) {
    case []byte: str = string(value.([]byte))
    case string: str = value.(string)
    default:
        return ErrInvalidType
    }

    s := strings.Split(str, ":")

    if len(s) != 3 {
        return ErrInvalidType
    }

    i, err := strconv.ParseInt(s[0], 10, 8)

    if err != nil {
        return err
    }

    d.Minutes = int(i * 60)

    i, err = strconv.ParseInt(s[1], 10, 8)

    if err != nil {
        return err
    }

    d.Minutes += int(i)

    return nil
}

func (d *TimeMin) MarshalJSON() ([]byte, error) {
    return json.Marshal(d.Minutes)
}

func (d *TimeMin) UnmarshalJSON(b []byte) error {
    return json.Unmarshal(b, &d.Minutes)
}

func (d *TimeMin) Val() int {
    return d.Minutes
}

func (d TimeMin) String() string {
    var (
        buf [32]byte
        w   int
    )

    w = len(buf)
    w = utils.UintToBStrLeadZero(buf[:w], 0)
    w--
    buf[w] = ':'
    w = utils.UintToBStrLeadZero(buf[:w], uint64(d.Minutes % 60))
    w--
    buf[w] = ':'
    w = utils.UintToBStrLeadZero(buf[:w], uint64(d.Minutes / 60))

    return string(buf[w:])
}

type TimeEmpty struct {
    Valid bool
    Time time.Time
}

func (t TimeEmpty) Value() (driver.Value, error) {
    if !t.Valid {
        return nil, nil
    }
    return t.String(), nil
}

func (t *TimeEmpty) Scan(value interface{}) error {
    if value == nil {
        t.Time, t.Valid = time.Time{}, false
        return nil
    }

    switch v := value.(type) {
    case time.Time: t.Time = v
    case []byte:    t.Time = utils.ParseDateTime(v, time.UTC)
    case string:    t.Time = utils.ParseDateTime([]byte(v), time.UTC)
    default:
        t.Valid = false
        return ErrInvalidType
    }

    t.Valid = true
    return nil
}

func (t TimeEmpty) String() string {
    var (
        buf [19]byte
        w   int
    )
 
    if !t.Valid {
        return "-"
    }

    w = len(buf)
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Second()))
    w--
    buf[w] = ':'
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Minute()))
    w--
    buf[w] = ':'
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Hour()))
    w--
    buf[w] = ' '
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Day()))
    w--
    buf[w] = '-'
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Month()))
    w--
    buf[w] = '-'
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Year()))

    return string(buf[:])
}

func (t *TimeEmpty) MarshalJSON() ([]byte, error) {
    var (
        buf [16]byte
        w   int
    )

    if !t.Valid {
        return []byte("\"-\""), nil
    }

    w = len(buf)
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Year()))
    w--
    buf[w] = '.'
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Month()))
    w--
    buf[w] = '.'
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Day()))
    w--
    buf[w] = ' '
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Minute()))
    w--
    buf[w] = ':'
    w = utils.UintToBStrLeadZero(buf[:w], uint64(t.Time.Hour()))

    return []byte("\"" + string(buf[:]) + "\""), nil
}

func (t *TimeEmpty) UnmarshalJSON(b []byte) (err error) {
    ptime, err := time.Parse("15:04 02.01.2006", string(b))
    t.Valid = err == nil
    if t.Valid {
        t.Time  = ptime
    } else {
        t.Time = time.Time{}
    }
    return nil
}