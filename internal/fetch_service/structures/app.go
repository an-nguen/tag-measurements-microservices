package structures

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"tag-measurements-microservices/internal/fetch_service/api"
	"tag-measurements-microservices/pkg/dto"
	"tag-measurements-microservices/pkg/models"
)

type App struct {
	Config              FetchConfig
	WirelessTagAccounts []models.WirelessTagAccount
	CloudClients        []api.CloudClient
	DataDb              *gorm.DB
}

func (app *App) InitCloudClients() {
	if err := app.DataDb.Find(&app.WirelessTagAccounts).Error; err != nil {
		log.Error("Failed to get wireless tag accounts.")
		return
	}
	if len(app.WirelessTagAccounts) == 0 {
		panic("No wireless tag accounts in database!")
	}

	tmpClient := &http.Client{}
	// Get session ids
	for i, account := range app.WirelessTagAccounts {
		sessionId, err := api.GetSessionId(&dto.LoginDataRequest{
			Email:    account.Email,
			Password: account.Password,
		})
		if err != nil {
			log.Error("Failed to get session id.")
			time.Sleep(10 * time.Minute)
			app.InitCloudClients()
			return
		}
		tagManagers, err := api.GetTagManagersApi(sessionId, "http://wirelesstag.net", tmpClient)
		if err != nil {
			log.Error("Failed to get tag managers.")
			time.Sleep(10 * time.Minute)
			app.InitCloudClients()
			return
		}

		for _, tm := range tagManagers {
			wstClient := api.CloudClient{
				Client:             &http.Client{},
				HostUrl:            "http://my.wirelesstag.net",
				SessionId:          "WTAG=" + sessionId,
				WirelessTagAccount: account,
				DataDB:             app.DataDb,
			}
			tm.Email = account.Email
			wstClient.TagManager = tm
			err = wstClient.SelectTagManager(tm.Mac)
			if err != nil {
				log.Error("Failed to select tag manager.")
				time.Sleep(10 * time.Minute)
				app.InitCloudClients()
				return
			}

			app.CloudClients = append(app.CloudClients, wstClient)
			app.WirelessTagAccounts[i].TagManagers = append(app.WirelessTagAccounts[i].TagManagers, tm)
		}
	}
	log.Info("Successful.")
}

func (app *App) StoreTagManagers() {
	for _, account := range app.WirelessTagAccounts {
		if account.TagManagers != nil {
			if len(account.TagManagers) > 0 {
				var databaseTagManagers models.TagManagerList
				app.DataDb.Where("email = ?", account.Email).Find(&databaseTagManagers)
				if len(databaseTagManagers) > 0 {
					diff := account.TagManagers.Difference(databaseTagManagers, func(a models.TagManager, b models.TagManager) bool {
						b.Mac = strings.Replace(b.Mac, ":", "", -1)
						return a.EqualsMac(b)
					})
					intersect := account.TagManagers.Intersect(databaseTagManagers, func(a models.TagManager, b models.TagManager) bool {
						b.Mac = strings.Replace(b.Mac, ":", "", -1)
						return a.EqualsMac(b)
					})
					if len(diff) > 0 {
						for _, tm := range diff {
							app.DataDb.Create(&tm)
						}
					}
					if len(intersect) > 0 {
						for _, tm := range intersect {
							app.DataDb.Save(&tm)
						}
					}
				}
			}
		}
	}
}

// Application config struct
type FetchConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	NDays    int    `json:"n_days"`
}

/*		Public function: ReadAppConfig
 *		---------------------------
 *      Read json data from a filename in a program filepath and fill FetchConfig struct
 *
 *		returns: filled FetchConfig struct
 */
func (app *App) ReadAppConfig(filename string) {
	// Get full path to config
	filepath, _ := os.Getwd()
	filepath = strings.ReplaceAll(filepath, "\\", "/")
	file, fileErr := os.Open(filepath + filename)
	if fileErr != nil {
		panic(fileErr)
	}
	buffer, readErr := ioutil.ReadAll(file)

	if readErr != nil {
		panic("Read error!")
	}
	_ = file.Close()

	var config FetchConfig
	unmarshalErr := json.Unmarshal(buffer, &config)

	if unmarshalErr != nil {
		panic(unmarshalErr)
	}
	app.Config = config
}
