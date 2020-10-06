package sync_functions

import (
	"errors"
	"tag-measurements-microservices/internal/fetch_service/api"
	"tag-measurements-microservices/internal/fetch_service/fetch_functions"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/utils"
)

/*
 *		function: syncManagers
 *		----------------------------
 *		Synchronize cloud tag managers with database managers
 *
 *		returns synchronized tag managers
 */
func SyncManagers(client api.WstClient) error {
	cloudManagers := fetch_functions.FetchManagers(client)
	if len(cloudManagers) == 0 {
		return errors.New("tag managers is empty")
	}
	var databaseTagManagers []models.TagManager
	client.Connection.Where("email = ?", client.Email).Find(&databaseTagManagers)
	var newTagManagers []models.TagManager

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
				if err := client.Connection.Save(&val).Error; err != nil {
					utils.LogError("SyncManagers", err)
				}
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

	sqlInsertStmt := "insert into tag_manager (mac, name, email) values ($1, $2, $3)"
	if len(newTagManagers) > 0 {
		for _, mgr := range newTagManagers {
			_, err := client.Connection.DB().Exec(sqlInsertStmt, mgr.Mac, mgr.Name, mgr.Email)
			if err != nil {
				utils.LogError("storeMeasurement:handleResponse", err)
			}
		}
	}

	return nil
}
