package api

import (
	"errors"
	"net/http"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"tag-measurements-microservices/internal/fetch_service/types"
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
	NewMeasurements    []models.Measurement
	ModMeasurements    []models.Measurement
}

func (c CloudClient) GetTags() []models.Tag {
	response, err := c.GetTagListApi()
	if err != nil {
		log.Error(err)
		return nil
	}
	tags := response.Tags(c.TagManager.Mac)
	sort.SliceStable(tags, func(i, j int) bool {
		return tags[i].SlaveId < tags[j].SlaveId
	})
	for i, _ := range tags {
		tags[i].Mac = GetMacWithSemicolons(tags[i].Mac)
	}

	return tags
}

func (c *CloudClient) GetTagManagers() []models.TagManager {
	response, err := c.GetTagManagersApi()
	if err != nil {
		log.Error(err)
		return nil
	}
	tagManagers := response.TagManagers(c.WirelessTagAccount.Email)
	return tagManagers
}

func (c *CloudClient) Store(db *gorm.DB) {
	if c.NewMeasurements != nil {
		if len(c.NewMeasurements) > 0 {
			if err := db.Create(&c.NewMeasurements).Error; err != nil {
				log.Error(err)
			}
		}
	}
	if c.ModMeasurements != nil {
		if len(c.ModMeasurements) > 0 {
			if err := db.Model(&models.Measurement{}).Updates(&c.ModMeasurements).Error; err != nil {
				log.Error(err)
			}
		}
	}
}

func (c *CloudClient) HandleResponse(json dto.MultiTagStatsRawResponse, uuid string, dataType types.MeasurementDataType) error {
	start := time.Now()

	for _, obj := range json.D.Stats {
		for j, tod := range obj.Tods[0] {

			value := obj.Values[0][j]
			date, err := time.Parse("1/2/2006", obj.Date)
			if err != nil {
				log.Error("storeMeasurement:handleResponse", "Failed to parse time ", err)
				continue
			}
			date = date.Add(time.Second * time.Duration(tod))
			if value == 0 {
				continue
			}
			record := findMeasurement(date, uuid, c.ModMeasurements)
			if record == nil {
				record = findMeasurement(date, uuid, c.NewMeasurements)
				if record == nil {
					var dbRecord models.Measurement
					db := c.DataDB.Order("date desc").Where("date = ? and tag_uuid = ?", date, uuid).First(&dbRecord)
					if db.Error != nil {
						record = nil
					} else {
						record = &dbRecord
					}
				}
			}
			if record == nil {
				switch dataType {
				case types.Temperature:
					c.NewMeasurements = append(c.NewMeasurements, models.Measurement{
						Date:        date,
						Temperature: value,
						TagUUID:     uuid,
					})
				case types.Humidity:
					c.NewMeasurements = append(c.NewMeasurements, models.Measurement{
						Date:     date,
						Humidity: value,
						TagUUID:  uuid,
					})
				case types.Signal:
					c.NewMeasurements = append(c.NewMeasurements, models.Measurement{
						Date:    date,
						Signal:  value,
						TagUUID: uuid,
					})
				case types.BatteryVoltage:
					c.NewMeasurements = append(c.NewMeasurements, models.Measurement{
						Date:    date,
						Voltage: value,
						TagUUID: uuid,
					})
				}
			} else {
				switch dataType {
				case types.Temperature:
					if record.Temperature != value {
						record.Temperature = value
					}
				case types.Humidity:
					if record.Humidity != value {
						record.Humidity = value
					}
				case types.BatteryVoltage:
					if record.Voltage != value {
						record.Voltage = value
					}
				case types.Signal:
					if record.Signal != value {
						record.Signal = value
					}
				default:
					err = errors.New("unknown type")
				}
			}
			if err != nil {
				return err
			}
		}
	}

	elapsed := time.Since(start)
	log.Info("Done for " + elapsed.String())
	return nil
}

func (c *CloudClient) StoreTags(db *gorm.DB) {
	if len(c.Tags) > 0 {
		var databaseTags models.TagList
		mac := GetMacWithSemicolons(c.TagManager.Mac)
		if err := db.Where("mac = ?", mac).Find(&databaseTags).Error; err != nil {
			log.Error(err)
		}
		if len(databaseTags) > 0 {
			diff := c.Tags.Difference(databaseTags, func(a models.Tag, b models.Tag) bool {
				return a.Equals(b)
			})
			intersect := c.Tags.Intersect(databaseTags, func(a models.Tag, b models.Tag) bool {
				return a.Equals(b)
			})
			if len(diff) > 0 {
				for _, tag := range diff {
					db.Create(&tag)
				}
			}
			if len(intersect) > 0 {
				for _, tag := range intersect {
					db.Save(&tag)
				}
			}
		}
	}

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

func (c *CloudClient) GetMeasurements(dataType types.MeasurementDataType, n int) {
	if c.DataDB == nil {
		log.Error("Connection can't be nil")
		return
	}
	// Check data type
	if !(dataType == "temperature" || dataType == "cap" || dataType == "signal" || dataType == "batteryVolt") {
		log.Error("Unknown data type!")
		return
	}
	for _, tag := range c.Tags {
		var slaveIds []int
		date := time.Now().Add(-(time.Duration(n) * (24 * time.Hour)))

		slaveIds = append(slaveIds, tag.SlaveId)

		if len(slaveIds) > 0 {
			jsonResponse, err := c.GetMultiTagStatsRawApi(slaveIds, string(dataType), date, time.Now().Add(2*(24*time.Hour)))
			if err != nil {
				log.Error("failed to get raw stats -> fetch again next minute. error - ", err)
				return
			}

			if err := c.HandleResponse(jsonResponse, tag.UUID, dataType); err != nil {
				log.Error("Error: ", err.Error())
				return
			}
		}
	}
}

func findMeasurement(date time.Time, uuid string, measurements []models.Measurement) *models.Measurement {
	pos := sort.Search(len(measurements), func(n int) bool {
		return measurements[n].Date == date && measurements[n].TagUUID == uuid
	})
	if pos >= len(measurements) || measurements[pos].TagUUID != uuid || measurements[pos].Date != date {
		return nil
	} else {
		return &measurements[pos]
	}
}
