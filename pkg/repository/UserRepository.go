package repository

import (
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"Thermo-WH/pkg/models"
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

	if err := repo.DataSource.Find(&user, "username = $1", username).Error; err != nil {
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

	if repo.DataSource.
		Preload("Roles").
		Where("id = ?", user.ID).
		First(&user).RecordNotFound() {
		return "", models.User{}, errors.New("User not found")
	}

	return token, user, nil
}

func (repo UserRepository) UpdateUser(user models.User) (models.User, error) {
	var databaseUser models.User

	if err := repo.DataSource.Find(&databaseUser, "id = $1", databaseUser.ID).Error; err != nil {
		return models.User{}, err
	}

	hashedPassword, err := repo.hashPassword(user)
	if err != nil {
		return models.User{}, err
	}

	databaseUser.Username = user.Username
	databaseUser.Password = string(hashedPassword)
	databaseUser.Roles = user.Roles
	repo.DataSource.Save(&databaseUser)
	repo.DataSource.Model(&databaseUser).Association("Roles").Replace(user.Roles)

	return databaseUser, nil
}

func (repo UserRepository) DeleteUser(id int) error {
	if err := repo.DataSource.Delete(&models.User{}, id).Error; err != nil {
		return err
	} else {
		return nil
	}
}
