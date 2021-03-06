package repository

import (
	"errors"
	"gorm.io/gorm"

	"tag-measurements-microservices/pkg/models"
)

type WarehouseGroupRepository struct {
	DataSource *gorm.DB
}

func (r WarehouseGroupRepository) GetTemperatureZones() ([]models.TemperatureZone, error) {
	var groups []models.TemperatureZone
	r.DataSource.Find(&groups)
	if len(groups) == 0 {
		return []models.TemperatureZone{}, errors.New("failed to find warehouse groups")
	}

	return groups, nil
}

func (r WarehouseGroupRepository) CreateTemperatureZone(group models.TemperatureZone) (models.TemperatureZone, error) {
	group.ID = 0
	if err := r.DataSource.Create(&group).Error; err != nil {
		return models.TemperatureZone{}, err
	}

	return group, nil
}

func (r WarehouseGroupRepository) UpdateTemperatureZone(group models.TemperatureZone) (models.TemperatureZone, error) {
	var groupDb models.TemperatureZone
	r.DataSource.First(&groupDb, "id = ?", group.ID)
	r.DataSource.Model(&groupDb).Updates(models.TemperatureZone{
		Name:            group.Name,
		Description:     group.Description,
		LowerTempLimit:  group.LowerTempLimit,
		HigherTempLimit: group.HigherTempLimit,
		NotifyEmails:    group.NotifyEmails,
	})
	r.DataSource.Model(&groupDb).Association("Tags").Clear()
	r.DataSource.Model(&groupDb).Association("Tags").Replace(&group.Tags)

	return group, nil
}

func (r WarehouseGroupRepository) GetWarehouseGroup(id string) (models.TemperatureZone, error) {
	var group models.TemperatureZone
	if err := r.DataSource.Where("id = ?", id).Preload("Tags").First(&group).Error; err != nil {
		return group, err
	}
	return group, nil
}

func (r WarehouseGroupRepository) DeleteWarehouseGroup(id string) error {
	var group models.TemperatureZone
	r.DataSource.Preload("Tags").First(&group, id)
	r.DataSource.Model(&group).Association("Tags").Clear()
	if err := r.DataSource.Delete(&group).Error; err != nil {
		return err
	}
	return nil
}
