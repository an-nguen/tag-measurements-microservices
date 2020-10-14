package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/repository"
)

type WirelessTagAccountController struct {
	Repository repository.WstAccountRepository
}

func (c WirelessTagAccountController) GetAccounts(ctx *gin.Context) {
	accounts, err := c.Repository.GetWirelessTagAccounts()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, accounts)
	}
}

func (c WirelessTagAccountController) AddAccount(ctx *gin.Context) {
	var jsonReq models.WirelessTagAccount
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := c.Repository.CreateWirelessTagAccount(jsonReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

func (c WirelessTagAccountController) UpdateAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	var jsonReq models.WirelessTagAccount
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := c.Repository.UpdateWirelessTagAccount(id, jsonReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

func (c WirelessTagAccountController) DeleteAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	var jsonReq models.WirelessTagAccount
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Repository.DeleteWirelessTagAccount(id)
}
