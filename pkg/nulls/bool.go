package nulls

import (
	"database/sql"
	"encoding/json"
)

type Bool struct {
	sql.NullBool
}

func NewBool(n interface{}) Bool {
	if n == nil {
		return Bool{
			sql.NullBool{
				Valid: false,
			},
		}
	}
	return Bool{
		sql.NullBool{
			Bool:  n.(bool),
			Valid: true,
		},
	}
}

func (n Bool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

func (n *Bool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		n.Valid = true
		n.Bool = *b
	} else {
		n.Valid = false
	}
	return nil
}
