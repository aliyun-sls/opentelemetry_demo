package model

import (
	"database/sql/driver"
	"encoding/json"
)

type PicPath []string

func (t *PicPath) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytesValue, _ := value.([]byte)
	if len(bytesValue) == 0 {
		return nil
	}
	return json.Unmarshal(bytesValue, t)
}

func (t PicPath) Value() (driver.Value, error) {
	return json.Marshal(t)
}
