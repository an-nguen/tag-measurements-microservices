package fetch_functions

import (
	"tag-measurements-microservices/internal/fetch_service/api"
	"tag-measurements-microservices/pkg/models"
)

/*
 *		function: fetchTag
 *		----------------------------
 *		Fetch from cloud server tags to collection and transfer to array of Tag
 *
 *		returns: array of Tag
 */
func FetchTags(client api.WstClient) []models.Tag {
	var tags []models.Tag
	jsonMap, _ := client.GetTagList()

	for _, item := range jsonMap.D {
		var tag models.Tag
		tag.Name = item.Name
		tag.Mac = client.TagManager.Mac
		tag.UUID = item.UUID
		tags = append(tags, tag)
	}

	return tags
}
