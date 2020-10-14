package dto

type MultiTagStatsRawResponse struct {
	D struct {
		Type  string `json:"__type"`
		Stats []struct {
			Date         string      `json:"date"`
			Ids          []int       `json:"ids"`
			Values       [][]float64 `json:"values"`
			ValuesBase64 interface{} `json:"values_base64"`
			Tods         [][]int     `json:"tods"`
			TodsBase64   interface{} `json:"tods_base64"`
		}
		TempUnit []interface{} `json:"temp_unit"`
		Ids      []int         `json:"ids"`
		Names    []string      `json:"names"`
		Discons  []interface{} `json:"discons"`
	} `json:"d"`
}
