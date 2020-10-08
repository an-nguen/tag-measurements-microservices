package sync_functions

import (
	"errors"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"tag-measurements-microservices/internal/fetch_service/api"
	"tag-measurements-microservices/internal/fetch_service/fetch_functions"
	"tag-measurements-microservices/internal/fetch_service/structures"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/utils"
)

/*
 *		function: syncTags
 *		----------------------------
 *		Synchronize cloud tags with database tags
 *
 */
func SyncTags(client api.WstClient, app *structures.App) error {
	cloudTags := fetch_functions.FetchTags(client)
	if len(cloudTags) == 0 {
		return errors.New("no cloudTags")
	}

	var databaseTags []models.Tag
	var newTags []interface{}
	app.DataDb.
		Where("mac = ?", client.TagManager.Mac).
		Find(&databaseTags)

	if len(cloudTags) > 0 && len(databaseTags) == 0 {
		for _, t := range cloudTags {
			newTags = append(newTags, t)
		}
	}

	if len(databaseTags) > 0 && len(cloudTags) > 0 {
		dbMapTags := make(map[string]models.Tag)
		for _, dt := range databaseTags {
			dbMapTags[dt.UUID] = dt
		}

		for _, t := range cloudTags {
			if val, ok := dbMapTags[t.UUID]; ok {
				val.Name = t.Name
				val.Mac = t.Mac
				if err := app.DataDb.Save(&val).Error; err != nil {
					utils.LogError("SyncTags", err)
				}
				for key := range dbMapTags {
					if key == t.UUID {
						delete(dbMapTags, key)
					}
				}
			} else {
				newTags = append(newTags, t)
			}
		}

		for _, val := range dbMapTags {
			app.DataDb.Delete(&val)
		}
	}

	if len(newTags) > 0 {
		if err := gormbulk.BulkInsert(app.DataDb, newTags,
			2500); err != nil {
			utils.LogError("SyncTags", err)
		}
	}

	return nil
}
