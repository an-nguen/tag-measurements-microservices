package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"tag-measurements-microservices/pkg/models"
	"time"

	"tag-measurements-microservices/pkg/repository"
)

type MeasurementController struct {
	Repository repository.MeasurementRepository
}

type GetTagDataMsg struct {
	UUIDList  []string `json:"uuid_list"`
	StartDate string   `json:"start_date"`
	EndDate   string   `json:"end_date"`
}

func (c MeasurementController) GetMeasurementsByUUID(ctx *gin.Context) {
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

func (c MeasurementController) GetMeasurementsCSVByUUID(ctx *gin.Context) {
	uuidListStr := ctx.Query("uuidList")
	uuidList := strings.Split(uuidListStr, ",")
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

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

	var measurements map[string][]models.Measurement
	measurements = make(map[string][]models.Measurement)

	for _, uuid := range uuidList {
		var tag models.Tag
		c.Repository.DataSource.Table("tag").Where("uuid = ?", uuid).First(&tag)
		var measurementsTemp []models.Measurement
		c.Repository.DataSource.Table("measurement").
			Where("tag_uuid = ? and date BETWEEN ? AND ?", uuid, startDate, endDate).
			Order("date DESC").
			Find(&measurementsTemp)
		measurements[tag.Name] = measurementsTemp
	}
	responseString := ""
	if err != nil {
		panic(err)
	}

	responseString += "Имя тега;Дата;Температура;Влажность;Напряжение;Вольтаж;\n"
	for k, v := range measurements {
		for _, data := range v {
			responseString += fmt.Sprintf("%s;%s;%f;%f;%f;%f;\n",
				k, data.Date.Format("02.01.2006 15:04:05-0700"),
				data.Temperature, data.Humidity, data.Voltage, data.Signal)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"csv": responseString})
}
