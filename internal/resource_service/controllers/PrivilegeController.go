package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	privilege, err := c.Repository.CreatePrivilege(jsonReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, privilege)
}

func (c PrivilegeController) UpdatePrivilege(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var jsonReq models.Privilege
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jsonReq.ID = uint(id)
	privilege, err := c.Repository.UpdatePrivilege(jsonReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, privilege)
}

func (c PrivilegeController) DeletePrivilege(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.Repository.DeletePrivilege(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "")
}
