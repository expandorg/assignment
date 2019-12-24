package nulls

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Time struct {
	Time  time.Time
	Valid bool
}

func (n *Time) Scan(value interface{}) error {
	if value == nil {
		n.Time = time.Time{}
		n.Valid = false
		return nil
	}
	timeValue, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	n.Time = timeValue
	n.Valid = true
	return nil
}

func (n Time) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func NewTime(n interface{}) Time {
	if n == nil {
		return Time{
			Valid: false,
		}
	}
	return Time{
		Time:  n.(time.Time),
		Valid: true,
	}
}

func (n Time) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *Time) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		n.Valid = true
		n.Time = *t
	} else {
		n.Valid = false
	}
	return nil
}
