package models

type WirelessTagAccount struct {
	Email    string `json:"email" gorm:"primary_key"`
	Password string `json:"password" gorm:"not_null"`
}

func (WirelessTagAccount) TableName() string {
	return "wireless_tag_account"
}
