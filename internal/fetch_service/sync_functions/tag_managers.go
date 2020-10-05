package sync_functions

import (
	"Thermo-WH/internal/fetch_service/api"
	"Thermo-WH/internal/fetch_service/fetch_functions"
	"Thermo-WH/pkg/models"
	"Thermo-WH/pkg/utils"
	"errors"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"sync"
)

/*
 *		function: syncManagers
 *		----------------------------
 *		Synchronize cloud tag managers with database managers
 *
 *		returns synchronized tag managers
 */
func SyncManagers(client api.WstClient) error {
	var lock sync.Mutex
	cloudManagers := fetch_functions.FetchManagers(client)
	if len(cloudManagers) == 0 {
		return errors.New("tag managers is empty")
	}
	var databaseTagManagers []models.TagManager
	client.Connection.Where("email = ?", client.Email).Find(&databaseTagManagers)
	var newTagManagers []interface{}

	if len(cloudManagers) > 0 && len(databaseTagManagers) == 0 {
		for _, mgr := range cloudManagers {
			newTagManagers = append(newTagManagers, mgr)
		}
	}
	if len(cloudManagers) > 0 && len(databaseTagManagers) > 0 {
		databaseMapTagManagers := make(map[string]models.TagManager)
		for _, tagManager := range databaseTagManagers {
			databaseMapTagManagers[tagManager.Mac] = tagManager
		}

		for _, cloudManager := range cloudManagers {
			if val, ok := databaseMapTagManagers[cloudManager.Mac]; ok {
				val.Name = cloudManager.Name
				val.Email = cloudManager.Email
				lock.Lock()
				if err := client.Connection.Save(&val).Error; err != nil {
					utils.LogError("SyncManagers", err)
				}
				lock.Unlock()
				for key := range databaseMapTagManagers {
					if key == cloudManager.Mac {
						delete(databaseMapTagManagers, key)
					}
				}
			} else {
				newTagManagers = append(newTagManagers, cloudManager)
			}
		}
	}

	lock.Lock()
	if len(newTagManagers) > 0 {
		if err := gormbulk.BulkInsert(client.Connection, newTagManagers, 2500); err != nil {
			utils.LogError("SyncFunctions", err)
		}
	}
	lock.Unlock()

	return nil
}
