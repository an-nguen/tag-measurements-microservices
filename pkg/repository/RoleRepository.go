package repository

import (
	"github.com/jinzhu/gorm"
	"tag-measurements-microservices/pkg/models"
)

type RoleRepository struct {
	DataSource *gorm.DB
}

func (repo RoleRepository) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := repo.DataSource.Preload("Users").
		Preload("Privileges").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (repo RoleRepository) CreateRole(role models.Role) (models.Role, error) {
	if err := repo.DataSource.Create(&role).Error; err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (repo RoleRepository) UpdateRole(role models.Role) (models.Role, error) {
	var databaseRole models.Role

	if err := repo.DataSource.Preload("Privileges").Preload("Privileges").Where("id = ?", role.ID).First(&databaseRole).Error; err != nil {
		return models.Role{}, err
	}

	databaseRole.Name = role.Name
	databaseRole.Description = role.Description
	databaseRole.Privileges = role.Privileges
	databaseRole.Users = role.Users
	if err := repo.DataSource.Save(&databaseRole).Error; err != nil {
		return models.Role{}, err
	}
	return databaseRole, nil
}

func (repo RoleRepository) DeleteRole(id int) error {
	if err := repo.DataSource.Delete(&models.Role{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
