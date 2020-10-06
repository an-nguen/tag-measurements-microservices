package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"tag-measurements-microservices/internal/fetch_service/store_functions"

	"tag-measurements-microservices/pkg/utils"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"tag-measurements-microservices/internal/fetch_service/api"
	"tag-measurements-microservices/internal/fetch_service/structures"
	utils2 "tag-measurements-microservices/internal/fetch_service/structures"
	"tag-measurements-microservices/internal/fetch_service/sync_functions"
	"tag-measurements-microservices/pkg/datasource"
	"tag-measurements-microservices/pkg/dto"
)

const FILENAME = "/configs/config_fetch.json"
const MeasurementTableName = "measurement"

var wg sync.WaitGroup

func main() {
	utils.LogPrintln("FETCH SERVICE", "Launched.")
	var app structures.App
	// Read config file from work directory
	app.Config = utils2.ReadAppConfig(FILENAME)
	utils.LogPrintln("main", fmt.Sprintf("The configuration file loaded from '%s'.", FILENAME))
	// Create connection to wst accounts database and temperature data database
	app.DataDb = datasource.InitDatabaseConnection(app.Config.Host, app.Config.Port,
		app.Config.User, app.Config.Password, app.Config.DbName)
	/** Main loop **/
	/*
	 *	Every 15 seconds fetch accounts from database and compare with accounts that fetched before.
	 */
	mainLoop(&app)
}

func mainLoop(app *structures.App) {
	for {
		initWstClients(app)
		utils.LogPrintln("mainLoop", "Begin to fetch measurements")
		for _, client := range app.WstClients {
			wg.Add(1)
			go clientFetch(app, client)
		}
		wg.Wait()
		utils.LogPrintln("mainLoop", "End of fetching measurements. Wait for next iteration")
		for _, client := range app.WstClients {
			err := client.Connection.Close()
			if err != nil {
				utils.LogError("clientFetch", err)
			}
		}

		time.Sleep(15 * time.Second)
	}
}

func initWstClients(app *structures.App) {
	app.WstAccounts = api.GetWstAccounts(app.DataDb)
	if len(app.WstAccounts) == 0 {
		panic("No wst account in database!")
	}

	tmpClient := &http.Client{}
	// Get session ids
	for _, acc := range app.WstAccounts {
		sessionId, err := api.GetSessionId(&dto.LoginDataRequest{
			Email:    acc.Email,
			Password: acc.Password,
		})
		if err != nil {
			utils.LogError("initWstClients", "Failed to get session id.")
			time.Sleep(3 * time.Minute)
			initWstClients(app)
			return
		}
		tagManagers, err := api.GetTagManagers(sessionId, "http://my.wirelesstag.net", tmpClient)
		if err != nil {
			utils.LogError("initWstClients", "Failed to get tag managers.")
			time.Sleep(3 * time.Minute)
			initWstClients(app)
			return
		}

		for _, tm := range tagManagers {
			wstClient := api.WstClient{
				Client:               &http.Client{},
				MeasurementTableName: MeasurementTableName,
				HostUrl:              "http://my.wirelesstag.net",
				SessionId:            "WTAG=" + sessionId,
				Email:                acc.Email,
				Password:             acc.Password,
			}
			tm.Email = acc.Email
			wstClient.TagManager = tm
			err = wstClient.SelectTagManager(tm.Mac)
			if err != nil {
				utils.LogError("initWstClients", "Failed to select tag manager.")
				time.Sleep(3 * time.Minute)
				initWstClients(app)
				return
			}

			wstClient.Connection = datasource.InitDatabaseConnection(app.Config.Host, app.Config.Port,
				app.Config.User, app.Config.Password, app.Config.DbName)
			app.WstClients = append(app.WstClients, wstClient)
		}
	}
	utils.LogPrintln("initWstClients", "Successful.")
}

func clientFetch(app *structures.App, wstClient api.WstClient) {
	defer wg.Done()
	/** Fetch managers and compare existed **/
	_ = sync_functions.SyncManagers(wstClient)

	/** Fetch tags with all information and compare with existed **/
	_ = sync_functions.SyncTags(wstClient)

	store_functions.StoreMeasurement(wstClient, app, store_functions.Temperature)
	store_functions.StoreMeasurement(wstClient, app, store_functions.Humidity)
	store_functions.StoreMeasurement(wstClient, app, store_functions.BatteryVoltage)
	store_functions.StoreMeasurement(wstClient, app, store_functions.Signal)
}
