package dto

type TagManagersResponse struct {
	D []struct {
		Type  string `json:"__type"`
		Users []struct {
			Name      string `json:"name"`
			Id        string `json:"id"`
			IsLimited bool   `json:"isLimited"`
		}
		Name               string `json:"name"`
		MAC                string `json:"mac"`
		LinkedToMac        string `json:"linkedToMac"`
		NotifyOfflineEmail string `json:"notifyOfflineEmail"`
		AllowMore          bool   `json:"allow_more"`
		Selected           bool   `json:"selected"`
		NotifyOffline      bool   `json:"notifyOffline"`
		NotifyOfflinePush  bool   `json:"notifyOfflinePush"`
		Online             bool   `json:"online"`
		WirelessConfig     struct {
			DataRate       int  `json:"dataRate"`
			ActiveInterval int  `json:"activeInterval"`
			Freq           int  `json:"Freq"`
			UseCRC16       bool `json:"useCRC16"`
			UseCRC32       bool `json:"useCRC32"`
			PSID           int  `json:"psid"`
		} `json:"wirelessConfig"`
		RadioId    string `json:"radioId"`
		Rev        int    `json:"rev"`
		DBID       int    `json:"dbid"`
		MStaticMAC string `json:"mStaticMac"`
	} `json:"d"`
}
