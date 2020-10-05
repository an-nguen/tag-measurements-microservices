package models

import (
	"strings"
	"time"
)

type Tag struct {
	UUID             string             `json:"uuid" gorm:"type:uuid;not_null;primary_key"`
	Name             string             `json:"name" gorm:"size:100;not_null"`
	Mac              string             `json:"mac"`
	VerificationDate time.Time          `json:"verification_date"`
	TemperatureZones []*TemperatureZone `json:"temperature_zones" gorm:"many2many:temperature_zone_tag;"`
}

func (src Tag) Equals(dst Tag) bool {
	if strings.Compare(src.UUID, dst.UUID) == 0 &&
		strings.Compare(src.Mac, dst.Mac) == 0 &&
		strings.Compare(src.Name, dst.Name) == 0 {
		return true
	} else {
		return false
	}
}

func (Tag) TableName() string {
	return "tag"
}
