package models

import "time"

type Measurement struct {
	Id          int64     `json:"id" gorm:"type:serial;primary_key"`
	Date        time.Time `json:"date" gorm:"type:timestamp"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Voltage     float64   `json:"voltage"`
	Signal      float64   `json:"signal"`
	TagUUID     string    `json:"tag_uuid" gorm:"type:uuid;not null"`
}

func (Measurement) TableName() string {
	return "measurement"
}
