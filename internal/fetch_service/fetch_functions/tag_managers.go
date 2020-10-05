package fetch_functions

import (
	"Thermo-WH/internal/fetch_service/api"
	"Thermo-WH/pkg/models"
	"bytes"
	"strings"
)

// Interfaces are named collections of method signatures.

/*
 *		function: fetchManagers
 *		----------------------------
 *		Fetch from cloud server tag managers to collection and transfer to array of TagManager
 *v
 *		returns: array of TagManager
 */
func FetchManagers(wstClient api.WstClient) []models.TagManager {
	jsonMap, _ := wstClient.GetTagManagers()
	var managers []models.TagManager
	for _, item := range jsonMap.D {
		var manager models.TagManager
		manager.Name = item.Name
		manager.Mac = strings.ToLower(item.MAC)
		var buf bytes.Buffer
		for i, ch := range manager.Mac {
			buf.WriteRune(ch)
			if (i+1)%2 == 0 && i != len(manager.Mac)-1 {
				buf.WriteRune(':')
			}
		}
		manager.Mac = buf.String()
		manager.Email = wstClient.Email
		managers = append(managers, manager)
	}
	return managers
}
