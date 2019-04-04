package md

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// JSONTime format json time field by myself
type Time struct {
	time.Time
}

const (
	Layout_YYYYMM         = "2006-01"
	Layout_YYYYMMDD       = "2006-01-02"
	Layout_YYYYMMDDHHIISS = "2006-01-02 15:04:05"
)

func NewTime() Time {
	return CreateTime(time.Now())
}
func NewTimePtr() *Time {
	return CreateTimePtr(time.Now())
}
func CreateTime(value interface{}) Time {
	if value == nil {
		return Time{}
	}
	if v, ok := value.(time.Time); ok {
		return Time{v}
	}
	if v, ok := value.(string); ok {
		layout := Layout_YYYYMMDDHHIISS
		data := []rune(v)
		if len(data) >= len(Layout_YYYYMMDDHHIISS) {
			layout = Layout_YYYYMMDDHHIISS
			data = data[:len(Layout_YYYYMMDDHHIISS)]
		} else if len(data) >= len(Layout_YYYYMMDD) {
			layout = Layout_YYYYMMDD
			data = data[:len(Layout_YYYYMMDD)]
		} else if len(data) >= len(Layout_YYYYMM) {
			layout = Layout_YYYYMM
			data = data[:len(Layout_YYYYMM)]
		}
		now, _ := time.ParseInLocation(layout, string(data), time.Local)
		return Time{now}
	}
	return Time{time.Now()}
}
func CreateTimePtr(value interface{}) *Time {
	t := CreateTime(value)
	return &t
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t Time) MarshalJSON() ([]byte, error) {
	if t.UnixNano() <= 0 || t.Unix() <= 0 || t.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}
func (t *Time) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		*t = Time{}
		return nil
	}
	layout := Layout_YYYYMMDDHHIISS
	data = []byte(strings.Replace(string(data), `"`, "", -1))
	if len(data) >= len(Layout_YYYYMMDDHHIISS) {
		layout = Layout_YYYYMMDDHHIISS
		data = data[:len(Layout_YYYYMMDDHHIISS)]
	} else if len(data) >= len(Layout_YYYYMMDD) {
		layout = Layout_YYYYMMDD
		data = data[:len(Layout_YYYYMMDD)]
	} else if len(data) >= len(Layout_YYYYMM) {
		layout = Layout_YYYYMM
		data = data[:len(Layout_YYYYMM)]
	}
	now, _ := time.ParseInLocation(layout, string(data), time.Local)
	if now.UnixNano() < 0 || now.Unix() <= 0 {
		*t = Time{}
	} else {
		*t = Time{now}
	}

	return nil
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *Time) Scan(v interface{}) error {
	if v == nil {
		*t = Time{}
		return nil
	}
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
