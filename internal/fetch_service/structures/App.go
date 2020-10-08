package structures

import (
	"github.com/lib/pq"
	"sync"
	"tag-measurements-microservices/internal/fetch_service/api"
	"tag-measurements-microservices/pkg/datasource"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/utils"
	"time"

	"github.com/jinzhu/gorm"
)

type App struct {
	Config             FetchConfig
	WstAccounts        []api.WirelessTagAccount
	WstClients         []api.WstClient
	DataDb             *gorm.DB
	InsertMeasurements []models.Measurement
	UpdateMeasurements []models.Measurement
}

func (app App) AddInsertMeasurement(measurement models.Measurement) {
	var mutex sync.Mutex
	mutex.Lock()
	app.InsertMeasurements = append(app.InsertMeasurements, measurement)
	mutex.Unlock()
}

func (app App) AddInsertMeasurements(measurements []models.Measurement) {
	var mutex sync.Mutex
	mutex.Lock()
	app.InsertMeasurements = append(app.InsertMeasurements, measurements...)
	mutex.Unlock()
}

func (app App) AddUpdateMeasurement(measurement models.Measurement) {
	var mutex sync.Mutex
	mutex.Lock()
	app.UpdateMeasurements = append(app.UpdateMeasurements, measurement)
	mutex.Unlock()
}
func (app App) AddUpdateMeasurements(measurements []models.Measurement) {
	var mutex sync.Mutex
	mutex.Lock()
	app.UpdateMeasurements = append(app.UpdateMeasurements, measurements...)
	mutex.Unlock()
}

func (app App) ClearInsertMeasurements() {
	var mutex sync.Mutex
	mutex.Lock()
	app.InsertMeasurements = app.InsertMeasurements[:0]
	mutex.Unlock()
}
func (app App) ClearUpdateMeasurements() {
	var mutex sync.Mutex
	mutex.Lock()
	app.InsertMeasurements = app.InsertMeasurements[:0]
	mutex.Unlock()
}

func (app App) StoreMeasurementDatabase() {
	start := time.Now()

	db := datasource.InitDatabaseConnection(app.Config.Host, app.Config.Port, app.Config.User, app.Config.Password, app.Config.DbName)
	sqlUpdateStmt := "UPDATE measurement SET temperature = $2, humidity = $3, voltage = $4, signal = $5 WHERE id = $1"

	if len(app.InsertMeasurements) > 0 {
		txn, err := db.DB().Begin()
		if err != nil {
			utils.LogError("app:StoreMeasurementDatabase", err)
		}
		stmt, _ := txn.Prepare(pq.CopyIn("measurement", "date", "temperature", "humidity", "voltage", "signal", "tag_uuid"))
		for _, measurement := range app.InsertMeasurements {
			_, err := stmt.Exec(measurement.Date, measurement.Temperature,
				measurement.Humidity, measurement.Voltage, measurement.Signal, measurement.TagUUID)
			if err != nil {
				utils.LogError("app:StoreMeasurementDatabase", err)
			}
		}
		_, err = stmt.Exec()
		if err != nil {
			utils.LogError("app:StoreMeasurementDatabase", err)
		}

		err = stmt.Close()
		if err != nil {
			utils.LogError("app:StoreMeasurementDatabase", err)
		}
		err = txn.Commit()
		if err != nil {
			utils.LogError("app:StoreMeasurementDatabase", err)
		}
	}
	if len(app.UpdateMeasurements) > 0 {
		for _, measurement := range app.UpdateMeasurements {
			_, err := db.DB().Exec(sqlUpdateStmt, measurement.Id, measurement.Temperature, measurement.Humidity, measurement.Voltage, measurement.Signal)
			if err != nil {
				utils.LogError("app:StoreMeasurementDatabase", err)
			}
		}
	}
	err := db.Close()
	if err != nil {
		utils.LogError("app:StoreMeasurementDatabase", err)
	}
	elapsed := time.Since(start)
	utils.LogPrintln("app:StoreMeasurementDatabase", "Done for "+elapsed.String())

}
