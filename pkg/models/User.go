package models

type User struct {
	ID       uint    `json:"id" gorm:"type:serial;primary_key"`
	Username string  `json:"username" gorm:"size:256;unique;not_null"`
	Password string  `json:"password" gorm:"size:256;not_null"`
	Roles    []*Role `json:"roles" gorm:"many2many:users_roles"`
}
