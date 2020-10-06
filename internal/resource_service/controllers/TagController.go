package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tag-measurements-microservices/pkg"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/repository"
	"tag-measurements-microservices/pkg/utils"
)

type TagController struct {
	Repository   repository.TagRepository
	ProxyService *pkg.ProxyService
}

func (c TagController) GetTags(ctx *gin.Context) {
	temperatureZoneId := ctx.Query("temperature_zone_id")
	if len(temperatureZoneId) == 0 {
		tags, err := c.Repository.GetTags()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, tags)
		} else {
			ctx.JSON(http.StatusOK, tags)
		}
	} else {
		tags, err := c.Repository.GetTagsByTemperatureZone(temperatureZoneId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, tags)
		} else {
			ctx.JSON(http.StatusOK, tags)
		}
	}
}

func (c TagController) UpdateTag(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	var err error = nil

	var jsonReq models.Tag
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.LogPrintln("UpdateTag", fmt.Sprintf("Attempt to update tag with uuid %d.", uuid))
	jsonReq.UUID = uuid
	jsonReq, err = c.Repository.UpdateTag(jsonReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jsonReq)
}

//
func (c TagController) GetLatestTagDetails(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.ProxyService.Tags)
}
