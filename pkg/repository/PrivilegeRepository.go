package repository

import (
	"gorm.io/gorm"

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
	if err := r.DataSource.Create(&privilege).Error; err != nil {
		return models.Privilege{}, err
	}
	r.DataSource.First(&privilege, "name = ?", privilege.Name)
	return privilege, nil
}

func (r PrivilegeRepository) UpdatePrivilege(privilege models.Privilege) (models.Privilege, error) {
	var databasePrivilege models.Privilege
	if err := r.DataSource.Where("id = ?", privilege.ID).First(&databasePrivilege).Error; err != nil {
		return models.Privilege{}, err
	}
	r.DataSource.Model(&databasePrivilege).Updates(models.Privilege{
		Name: privilege.Name,
	})
	r.DataSource.Model(&databasePrivilege).Association("Roles").Clear()
	r.DataSource.Model(&databasePrivilege).Association("Roles").Replace(&privilege.Roles)
	r.DataSource.First(&databasePrivilege, "id = ?", databasePrivilege.ID)
	return databasePrivilege, nil
}

func (r PrivilegeRepository) DeletePrivilege(id int) error {
	if err := r.DataSource.Where("id = ?", id).Delete(&models.Privilege{}).Error; err != nil {
		return err
	}
	return nil
}
