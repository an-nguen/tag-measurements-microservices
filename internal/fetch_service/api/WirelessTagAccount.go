package api

import (
	"github.com/jinzhu/gorm"
	"strings"
)

type WirelessTagAccount struct {
	Email    string
	Password string
}

func (WirelessTagAccount) TableName() string {
	return "wireless_tag_account"
}

func GetWstAccounts(db *gorm.DB) []WirelessTagAccount {
	var accounts []WirelessTagAccount
	if err := db.Find(&accounts).Error; err != nil {
		panic(err)
	}

	return accounts
}

func (acc WirelessTagAccount) Compare(account WirelessTagAccount) bool {
	return (strings.Compare(acc.Email, account.Email) == 0) &&
		(strings.Compare(acc.Password, account.Password) == 0)
}
