package controllers

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"

	"tag-measurements-microservices/pkg"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/repository"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  32768,
	WriteBufferSize: 32768,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	uuid := ctx.Param("id")
	var err error = nil

	var jsonReq models.Tag
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info(fmt.Sprintf("Attempt to update tag with uuid %d.", uuid))
	jsonReq, err = c.Repository.UpdateTag(jsonReq, uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jsonReq)
}

//
func (c TagController) GetLatestTagDetails(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Error("failed to set websocket upgrade: ", err)
		return
	}
	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		msgStr := string(msg)
		if strings.HasPrefix(msgStr, "/latest") {
			res, err := json.Marshal(c.ProxyService.Tags)
			if err != nil {
				log.Error("GetLatestTagDetails", "failed to marshal tags: ", err)
			}
			_ = ws.WriteMessage(messageType, res)
		}
	}
}
