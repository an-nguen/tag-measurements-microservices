package store_functions

import (
	"errors"
	"tag-measurements-microservices/pkg/datasource"
	"time"

	"tag-measurements-microservices/internal/fetch_service/api"
	"tag-measurements-microservices/internal/fetch_service/structures"
	"tag-measurements-microservices/pkg/dto"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/utils"
)

type MeasurementDataType string

const (
	Temperature    MeasurementDataType = "temperature"
	Humidity                           = "cap"
	BatteryVoltage                     = "batteryVolt"
	Signal                             = "signal"
)

func StoreMeasurement(client api.WstClient, app *structures.App, dataType MeasurementDataType) {
	if app.DataDb == nil {
		utils.LogError("StoreMeasurement", "Connection can't be nil")
		return
	}
	// Check data type
	if !(dataType == "temperature" || dataType == "cap" || dataType == "signal" || dataType == "batteryVolt") {
		utils.LogError("StoreMeasurement", "Unknown data type!")
		return
	}

	// Fetch tags
	jsonMap, _ := client.GetTagList()

	for _, item := range jsonMap.D {
		date := time.Now().Add(-(time.Duration(app.Config.NDays) * (24 * time.Hour)))

		var slaveIds []int
		slaveIds = append(slaveIds, item.SlaveId)
		jsonResponse, err := client.GetMultiTagStatsRaw(slaveIds, string(dataType), date, time.Now())
		if err != nil {
			utils.LogError("StoreMeasurement", "failed to get raw stats -> fetch again next minute. error - ", err)
			break
		}

		if err := handleResponse(jsonResponse, app, item.UUID, dataType); err != nil {
			utils.LogError("StoreMeasurement", "Error: ", err.Error())
			break
		}
	}
}

func handleResponse(json dto.MultiTagStatsRawResponse, app *structures.App, uuid string, dataType MeasurementDataType) error {
	start := time.Now()
	conn := datasource.InitDatabaseConnection(app.Config.Host, app.Config.Port,
		app.Config.User, app.Config.Password, app.Config.DbName)
	var insertRecords []models.Measurement
	var updateRecords []models.Measurement

	for _, obj := range json.D.Stats {
		for i, tod := range obj.Tods[0] {

			value := obj.Values[0][i]
			thatDay, err := time.Parse("1/2/2006", obj.Date)
			if err != nil {
				utils.LogError("storeMeasurement:handleResponse", "Failed to parse time ", err)
				continue
			}
			thatDay = thatDay.Add(time.Second * time.Duration(tod))
			if value == 0 {
				continue
			}
			var record models.Measurement
			db := conn.Order("date desc").Where("date = ? and tag_uuid = ?", thatDay, uuid).First(&record)
			switch dataType {
			case Temperature:
				if db.RecordNotFound() {
					insertRecords = append(insertRecords, models.Measurement{
						Date:        thatDay,
						Temperature: value,
						TagUUID:     uuid,
					})
				} else {
					if record.Temperature != value {
						record.Temperature = value
						updateRecords = append(updateRecords, record)
					}
				}
			case Humidity:
				if db.RecordNotFound() {
					insertRecords = append(insertRecords, models.Measurement{
						Date:     thatDay,
						Humidity: value,
						TagUUID:  uuid,
					})
				} else {
					if record.Humidity != value {
						record.Humidity = value
						updateRecords = append(updateRecords, record)
					}
				}
			case Signal:
				if db.RecordNotFound() {
					insertRecords = append(insertRecords, models.Measurement{
						Date:    thatDay,
						Signal:  value,
						TagUUID: uuid,
					})
				} else {
					if record.Signal != value {
						record.Signal = value
						updateRecords = append(updateRecords, record)
					}
				}
			case BatteryVoltage:
				if db.RecordNotFound() {
					insertRecords = append(insertRecords, models.Measurement{
						Date:    thatDay,
						Voltage: value,
						TagUUID: uuid,
					})
				} else {
					if record.Voltage != value {
						record.Voltage = value
						updateRecords = append(updateRecords, record)
					}
				}
			default:
				err = errors.New("unknown type")
			}
			if err != nil {
				return err
			}
		}
	}

	app.AddInsertMeasurements(insertRecords)
	app.AddUpdateMeasurements(updateRecords)

	elapsed := time.Since(start)
	conn.Close()
	utils.LogPrintln("storeMeasurement:handleResponse", "Done for "+elapsed.String())
	return nil
}
