package models

import (
	"strings"
)

type TagManager struct {
	Mac   string `json:"mac" gorm:"type:not_null;primary_key"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (src TagManager) EqualsAll(dst TagManager) bool {
	if strings.Compare(src.Name, dst.Name) == 0 &&
		strings.Compare(src.Mac, dst.Mac) == 0 &&
		strings.Compare(src.Email, dst.Email) == 0 {
		return true
	} else {
		return false
	}
}
func (src TagManager) EqualsMac(dst TagManager) bool {
	if strings.Compare(src.Mac, dst.Mac) == 0 {
		return true
	} else {
		return false
	}
}

func (TagManager) TableName() string {
	return "tag_manager"
}

type TagManagerList []TagManager

func (src TagManagerList) Difference(dst TagManagerList, compare func(a TagManager, b TagManager) bool) TagManagerList {
	var diff TagManagerList
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

func (src TagManagerList) Intersect(dst TagManagerList, compare func(a TagManager, b TagManager) bool) TagManagerList {
	var intersect TagManagerList
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
