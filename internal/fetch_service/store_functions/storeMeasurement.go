package store_functions

import (
	"errors"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"sync"
	"time"

	"Thermo-WH/internal/fetch_service/api"
	"Thermo-WH/internal/fetch_service/structures"
	"Thermo-WH/pkg/dto"
	"Thermo-WH/pkg/models"
	"Thermo-WH/pkg/utils"
)

/*
 *	parameter: dataType - it can be only: temperature/cap/signal/batteryVolt
 */
func StoreMeasurement(client api.WstClient, app *structures.App, dataType string) {
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
		if client.Connection.Find(&d, "tag_uuid = $1", item.UUID).Limit(1).RecordNotFound() {
			d.Date = time.Now().Add(-(time.Duration(app.Config.NDays) * (24 * time.Hour)))
		}

		jsonSpannedData, err := client.GetMultiTagStatsRaw(make([]int, item.SlaveId), dataType, d.Date, time.Now())
		for err != nil {
			utils.LogError("StoreMeasurement", "failed to get raw stats -> fetch again next minute. error - ", err)
			time.Sleep(1 * time.Minute)
			jsonSpannedData, err = client.GetMultiTagStatsRaw(make([]int, item.SlaveId), dataType, d.Date, time.Now())
		}

		if err := handleResponse(jsonSpannedData, client, item.UUID, dataType); err != nil {
			utils.LogError("StoreMeasurement", "Error: ", err.Error())
		}
	}
}

func handleResponse(jsonSpannedData dto.MultiTagStatsRawResponse, client api.WstClient, uuid string, t string) error {
	var insertRecords []interface{}
	var lock sync.Mutex

	for _, obj := range jsonSpannedData.D.Stats {
		for i, tod := range obj.Tods[0] {

			value := obj.Values[0][i]
			thatDay, err := time.Parse("1/2/2006", obj.Date)
			if err != nil {
				utils.LogError("handleResponse", "Failed to parse time ", err)
				continue
			}
			thatDay = thatDay.Add(time.Second * time.Duration(tod))
			if value == 0 {
				continue
			}
			var record models.Measurement
			lock.Lock()
			switch t {
			case "temperature":
				if client.Connection.First(&record,
					"date = ? and tag_uuid = ?", thatDay, uuid).RecordNotFound() {
					insertRecords = append(insertRecords, models.Measurement{
						Date:        thatDay,
						Temperature: value,
						TagUUID:     uuid,
					})
				} else {
					record.Temperature = value
					client.Connection.Save(&record)
				}
			case "cap":
				if client.Connection.First(&record,
					"date = ? and tag_uuid = ?", thatDay, uuid).RecordNotFound() {
					insertRecords = append(insertRecords, models.Measurement{
						Date:     thatDay,
						Humidity: value,
						TagUUID:  uuid,
					})
				} else {
					record.Humidity = value
					client.Connection.Save(&record)
				}
			case "signal":
				if client.Connection.First(&record,
					"date = ? and tag_uuid = ?", thatDay, uuid).RecordNotFound() {
					insertRecords = append(insertRecords, models.Measurement{
						Date:    thatDay,
						Signal:  value,
						TagUUID: uuid,
					})
				} else {
					record.Signal = value
					client.Connection.Save(&record)
				}
			case "batteryVolt":
				if client.Connection.First(&record,
					"date = ? and tag_uuid = ?", thatDay, uuid).RecordNotFound() {
					insertRecords = append(insertRecords, models.Measurement{
						Date:    thatDay,
						Voltage: value,
						TagUUID: uuid,
					})
				} else {
					record.Voltage = value
					client.Connection.Save(&record)
				}
			default:
				err = errors.New("unknown type")
			}
			lock.Unlock()
			if err != nil {
				return err
			}
		}
	}

	lock.Lock()
	if len(insertRecords) > 0 {
		if err := gormbulk.BulkInsert(client.Connection, insertRecords, 2500); err != nil {
			utils.LogError("handleResponse", err)
		}
	}
	lock.Unlock()
	return nil
}
