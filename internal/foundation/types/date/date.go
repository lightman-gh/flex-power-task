package date

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"
)

const Layout = "2006-01-02"

var _ sql.Scanner = (*ISO8601)(nil)
var _ driver.Valuer = (*ISO8601)(nil)

// ISO8601 is a wrapper for time.Time and supports
// to save and marshal ordinary ISO8601 date values.
type ISO8601 struct {
	time.Time
	nil bool
}

func NewISO8601Today() *ISO8601 {
	return NewISO8601(time.Now())
}

func NewISO8601(t time.Time) *ISO8601 {
	return &ISO8601{
		Time: t,
		nil:  t.IsZero(),
	}
}

func (d *ISO8601) ParseExample(v string) (interface{}, error) {
	x := NewISO8601Today()
	err := x.Scan(v)

	return x, err
}

func IsZero(d *ISO8601) bool {
	return d == nil || d.IsZero()
}

func (d *ISO8601) In(loc *time.Location) *ISO8601 {
	d.Time = d.Time.In(loc)

	return d
}

func (d *ISO8601) String() string {
	return d.Format(Layout)
}

func (d *ISO8601) Unwrap() *ISO8601 {
	if d == nil || d.nil {
		return nil
	}

	return d
}

func (d ISO8601) Value() (driver.Value, error) {
	if d.nil {
		return nil, nil
	}

	return []byte(d.String()), nil
}

func (d *ISO8601) Scan(src interface{}) error {
	d.nil = src == nil
	if d.nil {
		return nil
	}

	var err error

	switch cast := src.(type) {
	case time.Time:
		d.Time = cast
	case []byte:
		d.Time, err = time.Parse(Layout, string(cast))
	case string:
		d.Time, err = time.Parse(Layout, cast)
	}

	return err
}

func (d *ISO8601) MarshalJSON() ([]byte, error) {
	if d == nil || d.nil {
		return []byte("null"), nil
	}

	stamp := "\"" + d.String() + "\""

	return []byte(stamp), nil
}

func (d *ISO8601) UnmarshalJSON(data []byte) error {
	d.nil = string(data) == "null"
	if d.nil {
		return nil
	}

	stringDate := strings.Replace(string(data), "T00:00:00Z", "", 1)

	t, err := time.Parse(`"`+Layout+`"`, stringDate)
	d.Time = t

	return err
}

func (d *ISO8601) GetTimeAtMidnight() time.Time {
	return time.Date(d.Time.Year(), d.Time.Month(), d.Time.Day(), 0, 0, 0, 0, d.Time.Location())
}
