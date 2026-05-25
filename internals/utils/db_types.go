package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB struct {
	Data any `json:"data"`
}

func (j JSONB) Value() (driver.Value, error) {
	if j.Data == nil {
		return nil, nil
	}
	return json.Marshal(j.Data)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		j.Data = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &j.Data)
}
