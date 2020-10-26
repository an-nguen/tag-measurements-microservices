package api

import (
	"net/http"
	"sort"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"tag-measurements-microservices/pkg/dto"
	"tag-measurements-microservices/pkg/models"
)

type CloudClient struct {
	Client             *http.Client
	HostUrl            string
	SessionId          string
	TagManager         models.TagManager
	WirelessTagAccount models.WirelessTagAccount
	Tags               models.TagList
	DataDB             *gorm.DB
}

func (c CloudClient) GetTags() []dto.TagResponse {
	response, err := c.GetTagListApi()
	if err != nil {
		log.Error(err)
		return nil
	}
	tags := response.D
	sort.SliceStable(tags, func(i, j int) bool {
		return tags[i].SlaveId < tags[j].SlaveId
	})

	return tags
}

func GetMacWithSemicolons(src string) string {
	mac := ""
	for i, ch := range src {
		if (i+1)%2 == 0 && (i+1) != len(src) {
			mac += string(ch)
			mac += ":"
		} else {
			mac += string(ch)
		}
	}
	return mac
}

func (c *CloudClient) StoreRealTimeMeasurements(wg *sync.WaitGroup) {
	tags := c.GetTags()
	date := time.Now()
	var measurements []models.MeasurementRT
	for _, tag := range tags {
		measurement := models.MeasurementRT{
			Date:        date,
			Temperature: tag.Temperature,
			Humidity:    tag.Cap,
			Voltage:     tag.BatteryVolt,
			Signal:      float64(tag.SignaldBm),
			TagUUID:     tag.UUID,
		}
		measurements = append(measurements, measurement)
	}
	if len(measurements) > 0 {
		c.DataDB.Create(&measurements)
	}
	wg.Done()
}
