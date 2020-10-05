package models

type Role struct {
	ID          uint          `json:"id" gorm:"type:serial;primary_key"`
	Name        string        `json:"name" gorm:"size:256;unique"`
	Description string        `json:"description" gorm:"size:1024"`
	Users       []*User       `json:"users" gorm:"many2many:users_roles"`
	Permissions []*Permission `json:"permissions" gorm:"many2many:roles_permissions"`
}
