package models

type TemperatureZone struct {
	ID              uint    `json:"id" gorm:"type:serial;primary_key"`
	Name            string  `json:"name" gorm:"unique;not_null"`
	Description     string  `json:"description"`
	LowerTempLimit  float64 `json:"lower_temp_limit"`
	HigherTempLimit float64 `json:"higher_temp_limit"`
	NotifyEmails    string  `json:"notify_emails"`
	Tags            []Tag   `json:"tags" gorm:"many2many:temperature_zone_tags;"`
}

func (TemperatureZone) TableName() string {
	return "temperature_zone"
}
