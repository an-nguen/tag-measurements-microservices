package controllers

import (
	"Thermo-WH/pkg/repository"
	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	UserRepo repository.UserRepository
}

func (c PermissionController) GetPermissions(ctx *gin.Context) {

}
