package pkg

import (
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"sync"
	"tag-measurements-microservices/internal/fetch_service/api"
	structures2 "tag-measurements-microservices/internal/notify_service/structures"
	"tag-measurements-microservices/internal/resource_service/structures"
	"tag-measurements-microservices/pkg/datasource"
	"tag-measurements-microservices/pkg/dto"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/repository"
	"tag-measurements-microservices/pkg/utils"
	"time"
)

type TagDetail struct {
	UUID             string  `json:"uuid"`
	Comment          string  `json:"comment"`
	TagType          int     `json:"tagType"`
	LastComm         int64   `json:"lastComm"`
	Alive            bool    `json:"alive"`
	SignaldBm        int     `json:"signaldBm"`
	BatteryVolt      float64 `json:"batteryVolt"`
	Beeping          bool    `json:"beeping"`
	Lit              bool    `json:"lit"`
	OutOfRange       bool    `json:"outOfRange"`
	Lux              float64 `json:"lux"`
	Temperature      float64 `json:"temperature"`
	Cap              float64 `json:"cap"`
	CapRaw           float64 `json:"capRaw"`
	Txpwr            int64   `json:"txpwr"`
	BatteryRemaining float64 `json:"batteryRemaining"`
}

type ProxyService struct {
	Tags                []TagDetail
	WirelessTagAccounts []models.WirelessTagAccount
	WstClients          []api.WstClient
	DataDb              *gorm.DB
}

func (p *ProxyService) CreateProxyService(config structures.ResourceConfig) {
	p.DataDb = datasource.InitDatabaseConnection(config.Host, config.Port,
		config.User, config.Password, config.DbName)
	p.DataDb.Find(&p.WirelessTagAccounts)
	p.initWstClients()
}

func (p *ProxyService) CreateProxyService1(config structures2.NotifyConfig) {
	p.DataDb = datasource.InitDatabaseConnection(config.Host, config.Port,
		config.User, config.Password, config.DbName)
	p.DataDb.Find(&p.WirelessTagAccounts)
	p.initWstClients()
}

func (p *ProxyService) initWstClients() {
	tmpClient := &http.Client{}
	// Get session ids
	for _, acc := range p.WirelessTagAccounts {
		sessionId, err := api.GetSessionId(&dto.LoginDataRequest{
			Email:    acc.Email,
			Password: acc.Password,
		})
		if err != nil {
			utils.LogError("initWstClients", "Failed to get session id.")
			time.Sleep(3 * time.Minute)
			p.initWstClients()
			return
		}
		tagManagers, err := api.GetTagManagers(sessionId, "http://my.wirelesstag.net", tmpClient)
		if err != nil {
			utils.LogError("initWstClients", "Failed to get tag managers.")
			time.Sleep(3 * time.Minute)
			p.initWstClients()
			return
		}
		repo := repository.TagManagerRepository{DataSource: p.DataDb}

		for _, tm := range tagManagers {
			wstClient := api.WstClient{
				Client:    &http.Client{},
				HostUrl:   "http://my.wirelesstag.net",
				SessionId: "WTAG=" + sessionId,
				Email:     acc.Email,
				Password:  acc.Password,
			}
			dbTm, _ := repo.GetTagManagerByMac(tm.Mac)

			wstClient.TagManager = dbTm
			err = wstClient.SelectTagManager(tm.Mac)
			if err != nil {
				utils.LogError("initWstClients", "Failed to select tag manager.")
				time.Sleep(3 * time.Minute)
				p.initWstClients()
				return
			}

			p.WstClients = append(p.WstClients, wstClient)
		}
	}
}

func (p *ProxyService) Start() {
	var mutex sync.Mutex
	go func() {
		for {
			mutex.Lock()
			p.DataDb.Find(&p.WirelessTagAccounts)
			p.initWstClients()
			mutex.Unlock()

			p.Tags = p.Tags[:0]
			for _, client := range p.WstClients {
				res, err := client.GetTagList()
				if err != nil {
					log.Println(err)
				}
				for _, d := range res.D {
					var detail TagDetail
					detail.Alive = d.Alive
					detail.UUID = d.UUID
					detail.BatteryRemaining = d.BatteryRemaining
					detail.BatteryVolt = d.BatteryVolt
					detail.Beeping = d.Beeping
					detail.Temperature = d.Temperature
					detail.Cap = d.Cap
					detail.CapRaw = d.CapRaw
					detail.Lux = d.Lux
					detail.Lit = d.Lit
					detail.Comment = d.Comment
					detail.SignaldBm = d.SignaldBm
					detail.Txpwr = d.Txpwr
					p.Tags = append(p.Tags, detail)
				}
			}
			time.Sleep(15 * time.Second)
		}
	}()
}
