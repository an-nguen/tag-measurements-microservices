package repository

import (
	"Thermo-WH/pkg/models"
	"github.com/jinzhu/gorm"
)

type RoleRepository struct {
	DataSource *gorm.DB
}

func (repo RoleRepository) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := repo.DataSource.Find(&roles).Error; err != nil {
		return roles, err
	}
	return roles, nil
}

func (repo RoleRepository) CreateRole(role models.Role) (models.Role, error) {
	if err := repo.DataSource.Create(&role).Error; err != nil {
		return models.Role{}, err
	}
	return role, nil
}
