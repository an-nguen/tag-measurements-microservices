package models

import (
	"strings"
	"time"
)

type Tag struct {
	UUID                   string             `json:"uuid" gorm:"type:uuid;not_null;primary_key"`
	Name                   string             `json:"name" gorm:"size:100;not_null"`
	Mac                    string             `json:"mac"`
	VerificationDate       time.Time          `json:"verification_date"`
	HigherTemperatureLimit float64            `json:"higher_temperature_limit"`
	LowerTemperatureLimit  float64            `json:"lower_temperature_limit"`
	TemperatureZones       []*TemperatureZone `json:"temperature_zones" gorm:"many2many:temperature_zone_tags;"`
	SlaveId                int                `gorm:"-"`
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

type TagList []Tag

func (src TagList) Difference(dst TagList, compare func(a Tag, b Tag) bool) TagList {
	var diff TagList
	i := 0
	j := 0
	n := len(src)
	m := len(dst)
	for i < n && j < m {
		if !compare(src[i], dst[j]) {
			diff = append(diff, src[i])
		} else {
			j++
		}
		i++
	}
	for i < n {
		diff = append(diff, src[i])
	}
	return diff
}

func (src TagList) Intersect(dst TagList, compare func(a Tag, b Tag) bool) TagList {
	var intersect TagList
	i := 0
	j := 0
	n := len(src)
	m := len(dst)
	for i < n && j < m {
		if compare(src[i], dst[j]) {
			intersect = append(intersect, src[i])
		} else {
			j++
		}
		i++
	}
	for i < n {
		intersect = append(intersect, src[i])
	}
	return intersect
}
