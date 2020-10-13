package main

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"tag-measurements-microservices/internal/fetch_service/api"
	"tag-measurements-microservices/internal/fetch_service/structures"
	"tag-measurements-microservices/internal/fetch_service/types"
	"tag-measurements-microservices/pkg/datasource"
)

const FILENAME = "/configs/config_fetch.json"

func main() {
	log.Info("Launched.")
	var app structures.App
	// Read config file from work directory
	app.ReadAppConfig(FILENAME)
	log.Info(fmt.Sprintf("The configuration file loaded from '%s'.", FILENAME))
	/** Main loop **/
	/*
	 *	Every 15 seconds fetch accounts from database and compare with accounts that fetched before.
	 */
	mainLoop(&app)
}

func mainLoop(app *structures.App) {
	var wg sync.WaitGroup

	for {
		app.DataDb = datasource.InitDatabaseConnection(app.Config.Host, app.Config.Port,
			app.Config.User, app.Config.Password, app.Config.DbName)
		app.InitCloudClients()
		app.StoreTagManagers()
		log.Info("Begin to fetch measurements")
		for _, client := range app.CloudClients {
			clientFetch(app, client, &wg)
		}
		log.Info("Wait for goroutines")
		log.Info("End of fetching measurements. Wait for next iteration")

		time.Sleep(30 * time.Second)
	}
}

func clientFetch(app *structures.App, cloudClient api.CloudClient, wg *sync.WaitGroup) {
	/** Fetch tags with all information **/
	tags := cloudClient.GetTags()
	if tags == nil {
		return
	}
	cloudClient.Tags = tags
	cloudClient.StoreTags(app.DataDb)

	cloudClient.GetMeasurements(types.Temperature, app.Config.NDays)
	cloudClient.GetMeasurements(types.Humidity, app.Config.NDays)
	cloudClient.GetMeasurements(types.BatteryVoltage, app.Config.NDays)
	cloudClient.GetMeasurements(types.Signal, app.Config.NDays)
	cloudClient.Store(app.DataDb)
}
