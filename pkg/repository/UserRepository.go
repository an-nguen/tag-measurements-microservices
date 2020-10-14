package repository

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"tag-measurements-microservices/pkg/models"
)

type UserRepository struct {
	DataSource *gorm.DB
}

func (repo UserRepository) FindUsers() ([]models.User, error) {
	var users []models.User
	if err := repo.DataSource.Find(&users).Error; err != nil {
		return []models.User{}, err
	}

	return users, nil
}

func (repo UserRepository) CreateUser(user models.User) (models.User, error) {
	hashedPassword, err := repo.hashPassword(user)
	if err != nil {
		return models.User{}, err
	}

	user.ID = 0
	user.Password = string(hashedPassword)

	if err := repo.DataSource.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	repo.DataSource.Preload("Roles").First(&user)
	user.Password = ""

	return user, nil
}

func (repo UserRepository) hashPassword(user models.User) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, err
}

func (repo UserRepository) AuthUser(username string, password string, secret string) (string, models.User, error) {
	var user models.User

	if err := repo.DataSource.Find(&user, "username = ?", username).Error; err != nil {
		return "", models.User{}, err
	}

	userPass := []byte(password)
	hashedPass := []byte(user.Password)
	if err := bcrypt.CompareHashAndPassword(hashedPass, userPass); err != nil {
		return "", models.User{}, err
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Audience:  strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(120 * time.Minute).Unix(),
		Issuer:    "Authentication service",
		IssuedAt:  time.Now().Unix(),
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", models.User{}, err
	}

	if err := repo.DataSource.
		Preload("Roles").
		Where("id = ?", user.ID).
		First(&user).Error; err != nil {
		return "", models.User{}, err
	}

	return token, user, nil
}

func (repo UserRepository) UpdateUser(user models.User) (models.User, error) {
	var databaseUser models.User

	if err := repo.DataSource.First(&databaseUser, "id = ?", user.ID).Error; err != nil {
		return models.User{}, err
	}

	hashedPassword, err := repo.hashPassword(user)
	if err != nil {
		return models.User{}, err
	}
	repo.DataSource.Model(&databaseUser).Updates(models.User{
		Username: user.Username,
		Password: string(hashedPassword),
	})
	repo.DataSource.Model(&databaseUser).Association("Roles").Clear()
	repo.DataSource.Model(&databaseUser).Association("Roles").Replace(user.Roles)
	repo.DataSource.Preload("Roles").First(&databaseUser, "id = ?", databaseUser.ID)

	return databaseUser, nil
}

func (repo UserRepository) DeleteUser(id int) error {
	if err := repo.DataSource.Delete(&models.User{}, id).Error; err != nil {
		return err
	} else {
		return nil
	}
}

func (repo UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	if err := repo.DataSource.Select([]string{"id", "username"}).Preload("Roles").Preload("Roles.Privileges").Find(&users).Error; err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (repo UserRepository) GetUser(id int) (models.User, error) {
	var user models.User
	if err := repo.DataSource.Where("id = ?", id).Select([]string{"id", "username"}).Preload("Roles").Preload("Roles.Privileges").First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
