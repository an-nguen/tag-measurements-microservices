package dto

type StatsRawResponse struct {
	D []struct {
		Type        string        `json:"__type"`
		Date        string        `json:"date"`
		Tods        []int         `json:"tods"`
		TodsBase64  []interface{} `json:"tods_base64"`
		Temps       []float64     `json:"temps"`
		TempsBase64 []float64     `json:"temps_base64"`
		Caps        []float64     `json:"caps"`
		CapsBase64  []float64     `json:"caps_base64"`
	} `json:"d"`
}
