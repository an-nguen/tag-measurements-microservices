package dto

import "tag-measurements-microservices/pkg/models"

type TagResponse struct {
	Type                string      `json:"__type"`
	DBID                int         `json:"dbid"`
	NotificationJS      string      `json:"notificationJS"`
	Name                string      `json:"name"`
	UUID                string      `json:"uuid"`
	Comment             string      `json:"comment"`
	SlaveId             int         `json:"slaveId"`
	TagType             int         `json:"tagType"`
	Discon              interface{} `json:"discon"`
	LastComm            int64       `json:"lastComm"`
	Alive               bool        `json:"alive"`
	SignaldBm           int         `json:"signaldBm"`
	BatteryVolt         float64     `json:"batteryVolt"`
	Beeping             bool        `json:"beeping"`
	Lit                 bool        `json:"lit"`
	MigrationPending    bool        `json:"migrationPending"`
	BeepDurationDefault int64       `json:"beepDurationDefault"`
	EventState          float64     `json:"eventState"`
	TempEventDate       float64     `json:"tempEventDate"`
	OutOfRange          bool        `json:"outOfRange"`
	TempSpurTh          float64     `json:"tempSpurTh"`
	Lux                 float64     `json:"lux"`
	Temperature         float64     `json:"temperature"`
	TempCalOffset       float64     `json:"tempCalOffset"`
	CapCalOffset        float64     `json:"capCalOffset"`
	ImageMD5            interface{} `json:"image_md5"`
	Cap                 float64     `json:"cap"`
	CapRaw              float64     `json:"capRaw"`
	Az2                 float64     `json:"az2"`
	CapEventState       float64     `json:"capEventState"`
	LightEventState     float64     `json:"lightEventState"`
	Shorted             bool        `json:"shorted"`
	Zmod                interface{} `json:"zmod"`
	Thermostat          interface{} `json:"thermostat"`
	Playback            interface{} `json:"playback"`
	PostBackInterval    int64       `json:"postBackInterval"`
	Rev                 int64       `json:"rev"`
	Version1            int64       `json:"version1"`
	FreqOffset          int64       `json:"freqOffset"`
	FreqCalApplied      int64       `json:"freqCalApplied"`
	ReviveEvery         int64       `json:"reviveEvery"`
	OorGrace            int64       `json:"oorGrace"`
	TempBL              interface{} `json:"tempBL"`
	CapBL               interface{} `json:"capBL"`
	LuxBL               interface{} `json:"luxBL"`
	LBTh                float64     `json:"LBTh"`
	EnLBN               bool        `json:"enLBN"`
	Txpwr               int64       `json:"txpwr"`
	RssiMode            bool        `json:"rssiMode"`
	Ds18                bool        `json:"ds18"`
	V2flag              int64       `json:"v2flag"`
	BatteryRemaining    float64     `json:"batteryRemaining"`
}

type TagListResponse struct {
	D []TagResponse
}

func (r TagListResponse) Tags(mac string) []models.Tag {
	var tags []models.Tag
	for _, cloudTag := range r.D {
		var tag models.Tag
		tag.UUID = cloudTag.UUID
		tag.Name = cloudTag.Name
		tag.Mac = mac
		tag.SlaveId = cloudTag.SlaveId
		tags = append(tags, tag)
	}
	return tags
}
