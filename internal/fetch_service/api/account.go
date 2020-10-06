package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"tag-measurements-microservices/pkg/dto"
	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/utils"
)

/* 		Public function: GetSessionId
 * 		------------------------------
 * 		fetch session id by passing login credentials
 *
 *		returns: session string id
 */
func GetSessionId(data *dto.LoginDataRequest) (string, error) {
	client := &http.Client{}
	jsonData, parseErr := utils.Serialize(data)

	if parseErr != nil {
		log.Println("[ERROR - GetSessionId] Parse error - ", parseErr)
		return "", parseErr
	}

	request := utils.CreateRequest("POST", "http://wirelesstag.net/ethAccount.asmx/SignIn",
		bytes.NewBuffer(jsonData))

	response, _, err := utils.SendRequest(client, request)
	if err != nil {
		utils.LogPrintln("GetSessionId", fmt.Sprintln("SendRequest error", err))
		return "", errors.New("GetSessionId: SendRequest - network error")
	}

	if response.StatusCode == 500 {
		return "", errors.New("GetSessionId: 500 internal error on remote")
	}
	if len(response.Cookies()) == 0 {
		return "", errors.New("GetSessionId: failed to get cookie")
	}

	return response.Cookies()[0].Value, nil
}

/*		Public function: SelectTagManager
 *		---------------------------
 *      Select tag manager by mac and return another cookie
 *
 *		returns: error
 */
func (wc *WstClient) SelectTagManager(MACAddr string) error {
	reqBody, _ := json.Marshal(map[string]string{
		"mac": MACAddr,
	})
	req := utils.CreateRequest("POST", wc.HostUrl+"/ethAccount.asmx/SelectTagManager",
		bytes.NewBuffer(reqBody))
	defer req.Body.Close()
	req.Header.Set("Cookie", wc.SessionId)
	resp, _, err := utils.SendRequest(wc.Client, req)
	if err != nil {
		log.Println("SelectTagManager", "Error from SendRequest")
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("error on cloud server")
	}

	if resp.Header.Get("Set-Cookie") != "" {
		wc.SessionId = resp.Header.Get("Set-Cookie")
		return nil
	}

	return errors.New("failed to set cookie")
}

/*		Public function: GetTagManagers
 *		---------------------------
 *      Fetch list of managers from remote API service
 *
 *		returns: dto.TagManagersResponse
 */
func (wc *WstClient) GetTagManagers() (dto.TagManagersResponse, error) {
	req := utils.CreateRequest("POST", wc.HostUrl+"/ethAccount.asmx/GetTagManagers",
		bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Cookie", wc.SessionId)
	_, body, err := utils.SendRequest(wc.Client, req)
	if err != nil {
		log.Println("GetTagManagers", "Error from SendRequest")
		return dto.TagManagersResponse{}, err
	}

	// Extract json to struct WstDto
	var jsonMap dto.TagManagersResponse
	_ = json.Unmarshal([]byte(body), &jsonMap)

	return jsonMap, nil
}

func GetTagManagers(sessionId string, hostUrl string, client *http.Client) ([]models.TagManager, error) {
	req := utils.CreateRequest("POST", hostUrl+"/ethAccount.asmx/GetTagManagers",
		bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Cookie", "WTAG="+sessionId)
	_, body, err := utils.SendRequest(client, req)
	if err != nil {
		utils.LogError("GetTagManagers", "Error from SendRequest")
		return nil, err
	}

	// Deserialize json to struct WstDtoParseMAC
	var jsonMap dto.WstDto
	err = json.Unmarshal([]byte(body), &jsonMap)
	if err != nil {
		utils.LogError("GetTagManagers", "Failed to deserialize json object to WstDtoParseMAC")
		return nil, err
	}

	var tagManagers []models.TagManager
	for _, tagManager := range jsonMap.D {
		macAddr := tagManager["mac"].(string)
		macAddr = strings.ToLower(macAddr)
		macAddr = strings.ReplaceAll(macAddr, ":", "")
		item := models.TagManager{
			Mac:  macAddr,
			Name: tagManager["name"].(string),
		}

		tagManagers = append(tagManagers, item)
	}

	return tagManagers, nil
}
