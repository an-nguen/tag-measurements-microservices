package main

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"tag-measurements-microservices/internal/fetch_service/structures"
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

	app.DataDb = datasource.InitDatabaseConnection(app.Config.Host, app.Config.Port,
		app.Config.User, app.Config.Password, app.Config.DbName)
	for {
		app.InitCloudClients()

		log.Info("Begin to fetch measurements")
		for _, client := range app.CloudClients {
			go client.StoreRealTimeMeasurements(&wg)
			wg.Add(1)
		}
		wg.Wait()
		log.Info("Wait for goroutines")
		log.Info("End of fetching measurements. Wait for next iteration")

		time.Sleep(60 * time.Second)
	}
}
