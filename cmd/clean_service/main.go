package main

import (
	"time"

	"tag-measurements-microservices/internal/clean_service/structures"
	"tag-measurements-microservices/pkg/datasource"
	"tag-measurements-microservices/pkg/models"
)

var FILENAME = "/configs/config_clean.json"

func main() {
	config := structures.ReadAppConfig(FILENAME)
	db := datasource.InitDatabaseConnection(config.Host, config.Port, config.User, config.Password, config.DbName)
	for {
		db.Where("date < now() - INTERVAL '75 DAYS'").Delete(&models.Measurement{})
		time.Sleep(time.Duration(config.CheckIntervalSec) * time.Second)
	}
}
