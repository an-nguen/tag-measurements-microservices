package models

import (
	"strings"
)

type TagManager struct {
	Mac   string `json:"mac" gorm:"type:macaddr;not_null;primary_key"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (src TagManager) Equals(dst TagManager) bool {
	if strings.Compare(src.Name, dst.Name) == 0 &&
		strings.Compare(src.Mac, dst.Mac) == 0 {
		return true
	} else {
		return false
	}
}

func (TagManager) TableName() string {
	return "tag_manager"
}
