package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Thermo-WH/pkg/repository"
)

type MeasurementController struct {
	Repository repository.MeasurementRepository
}

type GetTagDataMsg struct {
	UUIDList  []string `json:"uuid_list"`
	StartDate string   `json:"start_date"`
	EndDate   string   `json:"end_date"`
}

func (c MeasurementController) GetTempByUUID(ctx *gin.Context) {
	uuidListStr := ctx.Query("uuidList")
	uuidList := strings.Split(uuidListStr, ",")
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")
	epsilonStr := ctx.Query("epsilon")

	if len(uuidList) == 0 || len(startDateStr) == 0 || len(endDateStr) == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": errors.New("failed to parse URL queries").Error()})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	epsilon, err := strconv.ParseFloat(epsilonStr, 8)
	if err != nil {
		epsilon = 0.0
	}

	tagData, _ := c.Repository.GetMeasurementByUUIDs(uuidList, startDate, endDate, epsilon)
	ctx.JSON(http.StatusOK, tagData)
}
