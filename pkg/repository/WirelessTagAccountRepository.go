package repository

import (
	"errors"
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"tag-measurements-microservices/pkg/models"
)

type WstAccountRepository struct {
	DataSource *gorm.DB
}

func (repo WstAccountRepository) GetWirelessTagAccounts() ([]models.WirelessTagAccount, error) {
	var accounts []models.WirelessTagAccount
	repo.DataSource.Find(&accounts)
	if len(accounts) == 0 {
		return nil, errors.New("failed to find accounts")
	}

	return accounts, nil
}

func (repo WstAccountRepository) CreateWirelessTagAccount(account models.WirelessTagAccount) (models.WirelessTagAccount, error) {
	if err := repo.DataSource.Create(&account).Error; err != nil {
		return models.WirelessTagAccount{}, errors.New("failed to create record")
	}
	return account, nil
}

func (repo WstAccountRepository) UpdateWirelessTagAccount(id string, acc models.WirelessTagAccount) (models.WirelessTagAccount, error) {
	var accountDb models.WirelessTagAccount
	if err := repo.DataSource.Where("email = ?", id).First(&accountDb).Error; err != nil {
		return models.WirelessTagAccount{}, err
	}
	accountDb.Email = acc.Email
	accountDb.Password = acc.Password
	repo.DataSource.Save(&accountDb)

	return accountDb, nil
}

func (repo WstAccountRepository) DeleteWirelessTagAccount(id string) {
	var accountDb models.WirelessTagAccount
	if err := repo.DataSource.Where("email = ?", id).First(&accountDb).Error; err != nil {
		return
	}
	repo.DataSource.Delete(&accountDb)
}
