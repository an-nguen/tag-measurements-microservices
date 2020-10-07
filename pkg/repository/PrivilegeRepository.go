package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"tag-measurements-microservices/pkg/models"
)

type PrivilegeRepository struct {
	DataSource *gorm.DB
}

func (r PrivilegeRepository) GetPrivileges() ([]models.Privilege, error) {
	var privileges []models.Privilege

	if err := r.DataSource.Preload("Roles").Find(&privileges).Error; err != nil {
		return []models.Privilege{}, err
	}
	return privileges, nil
}

func (r PrivilegeRepository) CreatePrivilege(privilege models.Privilege) (models.Privilege, error) {
	if r.DataSource.Where("name = ?", privilege.Name).First(&models.Privilege{}).RecordNotFound() {
		if err := r.DataSource.Create(&privilege).Error; err != nil {
			return models.Privilege{}, err
		}

		return privilege, nil
	} else {
		return models.Privilege{}, errors.New("name should be unique")
	}
}

func (r PrivilegeRepository) UpdatePrivilege(privilege models.Privilege) (models.Privilege, error) {
	var databasePrivilege models.Privilege
	if err := r.DataSource.Preload("Roles").Where("id = ?", privilege.ID).First(&databasePrivilege).Error; err != nil {
		return models.Privilege{}, err
	}

	databasePrivilege.Name = privilege.Name
	databasePrivilege.Roles = privilege.Roles
	r.DataSource.Save(&databasePrivilege)
	r.DataSource.Model(&databasePrivilege).Association("Roles").Replace(privilege.Roles)

	return databasePrivilege, nil
}

func (r PrivilegeRepository) DeletePrivilege(id int) error {
	if err := r.DataSource.Where("id = ?", id).Delete(&models.Privilege{}).Error; err != nil {
		return err
	}
	return nil
}
