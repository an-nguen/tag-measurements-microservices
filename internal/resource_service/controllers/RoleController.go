package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/repository"
	"tag-measurements-microservices/pkg/utils"
)

type RoleController struct {
	Secret     string
	Repository repository.RoleRepository
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
	if err := c.Repository.DataSource.
		Preload("Roles").
		Preload("Roles.Privileges").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "no user"})
		return
	}

	if len(user.Roles) == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "no roles"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": user.Roles})
}

func (c RoleController) GetRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := c.Repository.GetRole(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (c RoleController) GetRoles(ctx *gin.Context) {
	roles, err := c.Repository.GetRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

func (c RoleController) CreateRole(ctx *gin.Context) {
	var jsonReq models.Role
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	role, err := c.Repository.CreateRole(jsonReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (c RoleController) UpdateRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var jsonReq models.Role
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jsonReq.ID = uint(id)
	role, err := c.Repository.UpdateRole(jsonReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (c RoleController) DeleteRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.Repository.DeleteRole(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "")
}
