package types

import (
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
