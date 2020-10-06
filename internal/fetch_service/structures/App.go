package structures

import (
	"sync"
	"tag-measurements-microservices/internal/fetch_service/api"

	"github.com/jinzhu/gorm"
)

type App struct {
	Config      FetchConfig
	WstAccounts []api.WirelessTagAccount
	WstClients  []api.WstClient
	DataDb      *gorm.DB
	WaitGroup   sync.WaitGroup
}
