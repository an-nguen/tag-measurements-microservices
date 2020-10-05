package main

import (
	"time"

	"Thermo-WH/internal/clean_service/structures"
	"Thermo-WH/pkg/datasource"
	"Thermo-WH/pkg/models"
)

var FILENAME = "/configs/config_clean.json"

func main() {
	config := structures.ReadAppConfig(FILENAME)
	db := datasource.InitDatabaseConnection(config.Host, config.Port, config.User, config.Password, config.DbName)
	for {
		db.Where("date < now() - INTERVAL '40 DAYS'").Delete(models.Measurement{})
		time.Sleep(time.Duration(config.CheckIntervalSec) * time.Second)
	}
}
