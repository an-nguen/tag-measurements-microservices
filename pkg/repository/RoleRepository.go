package repository

import (
	"gorm.io/gorm"

	"tag-measurements-microservices/pkg/models"
)

type RoleRepository struct {
	DataSource *gorm.DB
}

func (repo RoleRepository) GetRole(id int) (models.Role, error) {
	var role models.Role
	if err := repo.DataSource.Preload("Privileges").Find(&role).Error; err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (repo RoleRepository) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := repo.DataSource.
		Preload("Privileges").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (repo RoleRepository) CreateRole(role models.Role) (models.Role, error) {
	role.ID = 0
	if err := repo.DataSource.Create(&role).Error; err != nil {
		return models.Role{}, err
	}
	repo.DataSource.Preload("Privileges").First(&role, "name = ? and description = ?", role.Name, role.Description)

	return role, nil
}

func (repo RoleRepository) UpdateRole(role models.Role) (models.Role, error) {
	var databaseRole models.Role

	if err := repo.DataSource.Where("id = ?", role.ID).First(&databaseRole).Error; err != nil {
		return models.Role{}, err
	}

	repo.DataSource.Model(&databaseRole).Updates(models.Role{
		Name:        role.Name,
		Description: role.Description,
	})
	repo.DataSource.Model(&databaseRole).Association("Privileges").Clear()
	repo.DataSource.Model(&databaseRole).Association("Privileges").Replace(role.Privileges)
	repo.DataSource.Preload("Privileges").First(&databaseRole, "id = ?", databaseRole.ID)

	return databaseRole, nil
}

func (repo RoleRepository) DeleteRole(id int) error {
	if err := repo.DataSource.Delete(&models.Role{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
