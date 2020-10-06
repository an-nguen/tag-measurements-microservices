package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/repository"
)

type PrivilegeController struct {
	Repository repository.PrivilegeRepository
}

func (c PrivilegeController) GetPrivileges(ctx *gin.Context) {
	privileges, err := c.Repository.GetPrivileges()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, privileges)
}

func (c PrivilegeController) CreatePrivilege(ctx *gin.Context) {
	var jsonReq models.Privilege
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}
