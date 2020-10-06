package api

import (
	"bytes"
	"encoding/json"
	"time"

	"tag-measurements-microservices/pkg/dto"
	"tag-measurements-microservices/pkg/utils"
)

type ReqOpts struct {
	Id       string `json:"id"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
}

/**
  data types - temperature, cap, batteryVolt, signal
*/

func (wc *WstClient) GetMultiTagStatsRaw(ids []int, t string, fromDate time.Time, toDate time.Time) (dto.MultiTagStatsRawResponse, error) {
	if len(ids) <= 0 {
		return dto.MultiTagStatsRawResponse{}, nil
	}
	type Req struct {
		Ids      []int  `json:"ids"`
		FromDate string `json:"fromDate"`
		ToDate   string `json:"toDate"`
		Type     string `json:"type"`
	}
	reqOption := Req{
		Ids:      ids,
		FromDate: fromDate.Format("01/02/2006"),
		ToDate:   toDate.Format("01/02/2006"),
		Type:     t,
	}
	jsonOption, err := utils.Serialize(reqOption)
	if err != nil {
		return dto.MultiTagStatsRawResponse{}, err
	}

	req := utils.CreateRequest("POST", wc.HostUrl+"/ethLogs.asmx/GetMultiTagStatsRaw", bytes.NewBuffer(jsonOption))
	req.Header.Add("Cookie", wc.SessionId)
	_, body, err := utils.SendRequest(wc.Client, req)
	if err != nil {
		return dto.MultiTagStatsRawResponse{}, err
	}

	defer req.Body.Close()
	var jsonMap dto.MultiTagStatsRawResponse
	_ = json.Unmarshal([]byte(body), &jsonMap)

	return jsonMap, nil
}
