package api

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

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
	jsonData, err := utils.Serialize(data)

	if err != nil {
		log.Error("Parse error - ", err)
		return "", err
	}

	request := utils.CreateRequest("POST", "http://wirelesstag.net/ethAccount.asmx/SignIn",
		bytes.NewBuffer(jsonData))

	response, _, err := utils.SendRequest(client, request)
	if err != nil {
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
func (c *CloudClient) SelectTagManager(MACAddr string) error {
	reqBody, _ := json.Marshal(map[string]string{
		"mac": MACAddr,
	})
	req := utils.CreateRequest("POST", c.HostUrl+"/ethAccount.asmx/SelectTagManager",
		bytes.NewBuffer(reqBody))
	defer req.Body.Close()
	req.Header.Set("Cookie", c.SessionId)
	resp, _, err := utils.SendRequest(c.Client, req)
	if err != nil {
		log.Error(err)
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("error on cloud server")
	}

	if resp.Header.Get("Set-Cookie") != "" {
		c.SessionId = resp.Header.Get("Set-Cookie")
		return nil
	}

	return errors.New("failed to set cookie")
}

/*		Public function: GetTagManagersApi
 *		---------------------------
 *      Fetch list of managers from remote API service
 *
 *		returns: dto.TagManagersResponse
 */
func (c *CloudClient) GetTagManagersApi() (dto.TagManagersResponse, error) {
	req := utils.CreateRequest("POST", c.HostUrl+"/ethAccount.asmx/GetTagManagers",
		bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Cookie", "WTAG="+c.SessionId)
	_, body, err := utils.SendRequest(c.Client, req)
	if err != nil {
		log.Error(err)
		return dto.TagManagersResponse{}, err
	}

	// Extract json to struct WirelessTagResponse
	var jsonMap dto.TagManagersResponse
	_ = json.Unmarshal([]byte(body), &jsonMap)

	return jsonMap, nil
}

func GetTagManagersApi(sessionId string, hostUrl string, client *http.Client) ([]models.TagManager, error) {
	req := utils.CreateRequest("POST", hostUrl+"/ethAccount.asmx/GetTagManagers",
		bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Cookie", "WTAG="+sessionId)
	_, body, err := utils.SendRequest(client, req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Deserialize json to struct WstDtoParseMAC
	var jsonMap dto.WirelessTagResponse
	err = json.Unmarshal([]byte(body), &jsonMap)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var tagManagers []models.TagManager
	for _, item := range jsonMap.D {
		macAddr := item.Mac
		macAddr = strings.ToLower(macAddr)
		macAddr = strings.ReplaceAll(macAddr, ":", "")
		item := models.TagManager{
			Mac:  macAddr,
			Name: item.Name,
		}

		tagManagers = append(tagManagers, item)
	}

	return tagManagers, nil
}
