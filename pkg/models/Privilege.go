package models

type Privilege struct {
	ID    uint    `json:"id" gorm:"type:serial;primary_key"`
	Name  string  `json:"name" gorm:"size:256;unique"`
	Roles []*Role `json:"roles" gorm:"many2many:roles_privileges"`
}
