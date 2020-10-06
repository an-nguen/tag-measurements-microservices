package store_functions

import (
	"errors"
	"github.com/lib/pq"
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

/*
 *	parameter: dataType - it can be only: temperature/cap/signal/batteryVolt
 */
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
		var d models.Measurement
		if client.Connection.Where("tag_uuid = ?", item.UUID).Order("date desc").First(&d).RecordNotFound() {
			d.Date = time.Now().Add(-(time.Duration(app.Config.NDays) * (24 * time.Hour)))
		}

		var slaveIds []int
		slaveIds = append(slaveIds, item.SlaveId)
		jsonSpannedData, err := client.GetMultiTagStatsRaw(slaveIds, string(dataType), d.Date, time.Now())
		if err != nil {
			utils.LogError("StoreMeasurement", "failed to get raw stats -> fetch again next minute. error - ", err)
			break
		}

		if err := handleResponse(jsonSpannedData, client, item.UUID, dataType); err != nil {
			utils.LogError("StoreMeasurement", "Error: ", err.Error())
			break
		}
	}
}

func handleResponse(jsonSpannedData dto.MultiTagStatsRawResponse, client api.WstClient, uuid string, dataType MeasurementDataType) error {
	var insertRecords []models.Measurement
	var updateRecords []models.Measurement

	for _, obj := range jsonSpannedData.D.Stats {
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
			db := client.Connection.First(&record,
				"date = ? and tag_uuid = ?", thatDay, uuid)
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

	sqlUpdateStmt := "UPDATE measurement SET temperature = $2, humidity = $3, voltage = $4, signal = $5 WHERE id = $1"
	if len(insertRecords) > 0 {
		txn, err := client.Connection.DB().Begin()
		if err != nil {
			utils.LogError("storeMeasurement:handleResponse", err)
		}
		stmt, _ := txn.Prepare(pq.CopyIn("measurement", "date", "temperature", "humidity", "voltage", "signal", "tag_uuid"))
		for _, measurement := range insertRecords {
			_, err := stmt.Exec(measurement.Date, measurement.Temperature,
				measurement.Humidity, measurement.Voltage, measurement.Signal, measurement.TagUUID)
			if err != nil {
				utils.LogError("storeMeasurement:handleResponse", err)
			}
		}
		_, err = stmt.Exec()
		if err != nil {
			utils.LogError("storeMeasurement:handleResponse", err)
		}

		err = stmt.Close()
		if err != nil {
			utils.LogError("storeMeasurement:handleResponse", err)
		}
		err = txn.Commit()
		if err != nil {
			utils.LogError("storeMeasurement:handleResponse", err)
		}
	}
	if len(updateRecords) > 0 {
		for _, measurement := range updateRecords {
			_, err := client.Connection.DB().Exec(sqlUpdateStmt, measurement.Id, measurement.Temperature, measurement.Humidity, measurement.Voltage, measurement.Signal)
			if err != nil {
				utils.LogError("storeMeasurement:handleResponse", err)
			}
		}
	}
	return nil
}
