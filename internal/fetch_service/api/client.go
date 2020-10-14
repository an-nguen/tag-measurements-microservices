package api

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"

	"tag-measurements-microservices/pkg/dto"
	"tag-measurements-microservices/pkg/utils"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

/*		Public function: GetTagListApi
 *		---------------------------
 *      Fetch list of temperature tags from remote API service
 *
 *		returns: dto.TagListResponse, error
 */
func (c *CloudClient) GetTagListApi() (dto.TagListResponse, error) {
	req := utils.CreateRequest("POST", c.HostUrl+"/ethClient.asmx/GetTagList",
		bytes.NewBuffer([]byte(`{}`)))
	req.Header.Add("Cookie", c.SessionId)
	_, body, err := utils.SendRequest(c.Client, req)
	if err != nil {
		return dto.TagListResponse{}, err
	}

	var jsonMap dto.TagListResponse
	_ = json.Unmarshal([]byte(body), &jsonMap)

	return jsonMap, err
}
