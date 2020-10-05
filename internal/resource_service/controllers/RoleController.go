package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"Thermo-WH/pkg/models"
	"Thermo-WH/pkg/repository"
	"Thermo-WH/pkg/utils"
)

type RoleController struct {
	Secret     string
	Repository repository.UserRepository
}

func (c RoleController) GetRolesByToken(ctx *gin.Context) {
	tokenStr, err := utils.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var claims jwt.StandardClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(c.Secret), nil
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if token.Valid == false {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	} else if claims.ExpiresAt < time.Now().Unix() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token is expired"})
		return
	}

	id, err := strconv.Atoi(claims.Audience)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if c.Repository.DataSource.
		Preload("Roles").
		Where("id = ?", id).
		First(&user).RecordNotFound() {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "no user"})
		return
	}

	if len(user.Roles) == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "no roles"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": user.Roles})
}
