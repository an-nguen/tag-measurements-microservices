package controllers

import (
	"Thermo-WH/pkg/models"
	"Thermo-WH/pkg/repository"
	"Thermo-WH/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TemperatureZoneController struct {
	Repository repository.WarehouseGroupRepository
}

func (c TemperatureZoneController) GetTemperatureZones(ctx *gin.Context) {
	groups, err := c.Repository.GetTemperatureZones()
	if err != nil {
		utils.LogError("GetTemperatureZones", err)
		return
	}
	ctx.JSON(http.StatusOK, groups)
}

func (c TemperatureZoneController) CreateTemperatureZone(ctx *gin.Context) {
	var jsonRes models.TemperatureZone
	if err := ctx.ShouldBindJSON(&jsonRes); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	retGroup, err := c.Repository.CreateTemperatureZone(jsonRes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, retGroup)
}

func (c TemperatureZoneController) UpdateTemperatureZone(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var jsonReq models.TemperatureZone
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonReq.ID = uint(id)
	jsonReq, err = c.Repository.UpdateTemperatureZone(jsonReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jsonReq)
}

func (c TemperatureZoneController) DeleteTemperatureZone(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.Repository.DeleteWarehouseGroup(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"id": id})

}

func (c TemperatureZoneController) GetTemperatureZone(ctx *gin.Context) {
	id := ctx.Param("id")
	group, err := c.Repository.GetWarehouseGroup(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, group)
}
