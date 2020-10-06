package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/repository"
	"tag-measurements-microservices/pkg/utils"
)

type TagManagerController struct {
	Repository repository.TagManagerRepository
}

func (c TagManagerController) GetTagManagers(ctx *gin.Context) {
	tagManagers, err := c.Repository.GetTagManagers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tagManagers)
}

func (c TagManagerController) UpdateTagManager(ctx *gin.Context) {
	mac := ctx.Param("mac")

	var jsonReq models.TagManager
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.LogPrintln("UpdateTagManager", fmt.Sprintf("Attempt to update tag manager with mac %d.", mac))
	jsonReq.Mac = mac
	jsonReq, err := c.Repository.UpdateTagManager(jsonReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jsonReq)
}

func (c TagManagerController) GetTagManager(ctx *gin.Context) {
	id := ctx.Param("mac")
	tagManager, err := c.Repository.GetTagManager(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tagManager)
}
