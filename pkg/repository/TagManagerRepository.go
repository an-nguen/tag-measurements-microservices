package repository

import (
	"gorm.io/gorm"

	"tag-measurements-microservices/internal/fetch_service/api"
	"tag-measurements-microservices/pkg/models"
)

type TagManagerRepository struct {
	DataSource *gorm.DB
}

func (r TagManagerRepository) GetTagManagers() ([]models.TagManager, error) {
	var tagManagers []models.TagManager
	if err := r.DataSource.Find(&tagManagers).Error; err != nil {
		return []models.TagManager{}, err
	}

	return tagManagers, nil
}

func (r TagManagerRepository) GetTagManagersByAccountName(accountName string) ([]models.TagManager, error) {
	var tagManagers []models.TagManager
	if err := r.DataSource.Find(&tagManagers, "account_name = ?", accountName).Error; err != nil {
		return tagManagers, err
	}

	return tagManagers, nil
}

func (r TagManagerRepository) UpdateTagManager(manager models.TagManager) (models.TagManager, error) {
	var managerDb models.TagManager
	r.DataSource.First(&managerDb, manager.Mac)
	r.DataSource.Save(&managerDb)

	return managerDb, nil
}

func (r TagManagerRepository) GetTagManagerByMac(mac string) (models.TagManager, error) {
	var manager models.TagManager
	if err := r.DataSource.First(&manager, "mac::varchar = ?", api.GetMacWithSemicolons(mac)).Error; err != nil {
		return models.TagManager{}, err
	}

	return manager, nil
}

func (r TagManagerRepository) GetTagManager(id string) (models.TagManager, error) {
	var tagManager models.TagManager
	if err := r.DataSource.Where("id = ?", id).Preload("TemperatureZones").First(&tagManager).Error; err != nil {
		return tagManager, err
	}

	return tagManager, nil
}
