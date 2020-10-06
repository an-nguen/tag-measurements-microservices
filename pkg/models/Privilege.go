package models

type Privilege struct {
	ID    uint    `json:"id" gorm:"type:serial;primary_key"`
	Name  string  `json:"name" gorm:"size:256;unique"`
	Value string  `json:"description" gorm:"size:256"`
	Roles []*Role `json:"roles" gorm:"many2many:roles_permissions"`
}
