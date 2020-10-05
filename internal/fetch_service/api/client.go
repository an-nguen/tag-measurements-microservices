package api

import (
	"bytes"
	"encoding/json"

	"Thermo-WH/pkg/dto"
	"Thermo-WH/pkg/utils"
)

/*		Public function: GetTagList
 *		---------------------------
 *      Fetch list of temperature tags from remote API service
 *
 *		returns: dto.TagListResponse, error
 */
func (wc *WstClient) GetTagList() (dto.TagListResponse, error) {
	req := utils.CreateRequest("POST", wc.HostUrl+"/ethClient.asmx/GetTagList",
		bytes.NewBuffer([]byte(`{}`)))
	req.Header.Add("Cookie", wc.SessionId)
	_, body, err := utils.SendRequest(wc.Client, req)
	if err != nil {
		return dto.TagListResponse{}, err
	}

	var jsonMap dto.TagListResponse
	_ = json.Unmarshal([]byte(body), &jsonMap)

	return jsonMap, err
}
