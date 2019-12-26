package nulls

import (
	"database/sql"
	"encoding/json"
)

type Int64 struct {
	sql.NullInt64
}

func NewInt64(n interface{}) Int64 {
	if n == nil {
		return Int64{
			sql.NullInt64{
				Valid: false,
			},
		}
	}
	return Int64{
		sql.NullInt64{
			Int64: n.(int64),
			Valid: true,
		},
	}
}

func (n *Int64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

func (n *Int64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		n.Valid = true
		n.Int64 = *i
	} else {
		n.Valid = false
	}
	return nil
}
