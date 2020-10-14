package dto

type WirelessTagResponse struct {
	D []struct {
		Type  string `json:"__type"`
		Users []struct {
			Name      string `json:"name"`
			Id        string `json:"id"`
			IsLimited bool   `json:"isLimited"`
		} `json:"users"`
		Name               string `json:"name"`
		Mac                string `json:"mac"`
		LinkedToMac        string `json:"linkedToMac"`
		NotifyOfflineEmail string `json:"notifyOfflineEmail"`
		AllowMore          bool   `json:"allowMore"`
		Selected           bool   `json:"selected"`
		NotifyOffline      bool   `json:"notifyOffline"`
		NotifyOfflinePush  bool   `json:"notifyOfflinePush"`
		Online             bool   `json:"online"`
		WirelessConfig     struct {
			DataRate       float64 `json:"dataRate"`
			ActiveInterval float64 `json:"activeInterval"`
			Freq           float64 `json:"Freq"`
			UseCRC16       bool    `json:"useCRC16"`
			UseCRC32       bool    `json:"useCRC32"`
			Psid           float64 `json:"psid"`
		}
		RadioId    string `json:"radioId"`
		Rev        int    `json:"rev"`
		DBID       int    `json:"dbid"`
		WsRoot     string `json:"wsRoot"`
		MStaticMac string `json:"mStaticMac"`
	} `json:"d"`
}
