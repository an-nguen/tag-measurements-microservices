package structures

import (
	"Thermo-WH/internal/fetch_service/api"
	"sync"

	"github.com/jinzhu/gorm"
)

type App struct {
	Config      FetchConfig
	WstAccounts []api.WirelessTagAccount
	WstClients  []api.WstClient
	DataDb      *gorm.DB
	WaitGroup   sync.WaitGroup
}
