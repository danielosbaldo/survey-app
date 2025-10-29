package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func (j JSONB) Value() (driver.Value, error) {
	b, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (j *JSONB) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for JSONB")
	}
	return json.Unmarshal(b, j)
}
