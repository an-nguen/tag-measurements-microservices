package api

import (
	"bytes"
	"time"

	"tag-measurements-microservices/pkg/dto"
	"tag-measurements-microservices/pkg/utils"
)

/**
  data types - temperature, cap, batteryVolt, signal
*/

func (c *CloudClient) GetMultiTagStatsRawApi(ids []int, t string, fromDate time.Time, toDate time.Time) (dto.MultiTagStatsRawResponse, error) {
	if len(ids) <= 0 {
		return dto.MultiTagStatsRawResponse{}, nil
	}
	type Request struct {
		Ids      []int  `json:"ids"`
		FromDate string `json:"fromDate"`
		ToDate   string `json:"toDate"`
		Type     string `json:"type"`
	}
	reqOption := Request{
		Ids:      ids,
		FromDate: fromDate.Format("01/02/2006"),
		ToDate:   toDate.Format("01/02/2006"),
		Type:     t,
	}
	jsonOption, err := utils.Serialize(reqOption)
	if err != nil {
		return dto.MultiTagStatsRawResponse{}, err
	}

	req := utils.CreateRequest("POST", c.HostUrl+"/ethLogs.asmx/GetMultiTagStatsRaw", bytes.NewBuffer(jsonOption))
	req.Header.Add("Cookie", c.SessionId)
	_, body, err := utils.SendRequest(c.Client, req)
	if err != nil {
		return dto.MultiTagStatsRawResponse{}, err
	}

	defer req.Body.Close()
	var jsonMap dto.MultiTagStatsRawResponse
	_ = json.Unmarshal([]byte(body), &jsonMap)

	return jsonMap, nil
}
