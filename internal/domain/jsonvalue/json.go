package jsonvalue

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSON []byte

func New(value string) (JSON, error) {
	if value == "" {
		value = "{}"
	}

	if !json.Valid([]byte(value)) {
		return nil, fmt.Errorf("metadata must be valid JSON")
	}

	return JSON(value), nil
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "{}", nil
	}

	return string(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON("{}")
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
	case string:
		*j = append((*j)[0:0], v...)
	default:
		return fmt.Errorf("unsupported JSON value type %T", value)
	}

	if !json.Valid(*j) {
		return fmt.Errorf("metadata must be valid JSON")
	}

	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("{}"), nil
	}

	return j, nil
}

func (j *JSON) UnmarshalJSON(value []byte) error {
	if !json.Valid(value) {
		return fmt.Errorf("metadata must be valid JSON")
	}

	*j = append((*j)[0:0], value...)
	return nil
}
