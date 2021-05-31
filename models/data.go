package models

import (
	"encoding/json"
	"errors"

	"database/sql/driver"
)

const (
	DateFormat  = "2006-01-02"
	TimeFormate = "2006-01-02 15:04:05"
)

// QDrama 自定义 Drama URL 结构
type QDrama struct {
	EP  string `json:"ep"`
	URL string `json:"url"`
}

type QDramaArray []QDrama

func (d *QDramaArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("no []byte")
	}
	return json.Unmarshal(b, &d)
}

func (d QDramaArray) Value() (driver.Value, error) {
	return json.Marshal(d)
}

type QMediaArray []string

func (m *QMediaArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &m)
}

func (m QMediaArray) Value() (driver.Value, error) {
	return json.Marshal(m)
}
