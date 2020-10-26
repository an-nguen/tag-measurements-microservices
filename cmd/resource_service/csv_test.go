package main

import (
	"fmt"
	"os"
	"tag-measurements-microservices/pkg/datasource"
	"testing"
	"time"
)

type Measurement struct {
	Id          int64     `json:"id" gorm:"type:serial;primary_key"`
	Date        time.Time `json:"date" gorm:"type:timestamp"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Voltage     float64   `json:"voltage"`
	Signal      float64   `json:"signal"`
	TagUUID     string    `json:"tag_uuid" gorm:"type:uuid;not null"`
}

type Tag struct {
	UUID             string    `json:"uuid" gorm:"type:uuid;not_null;primary_key"`
	Name             string    `json:"name" gorm:"size:100;not_null"`
	Mac              string    `json:"mac"`
	VerificationDate time.Time `json:"verification_date"`
}

func TestCreateCSV(t *testing.T) {
	var measurements map[string][]Measurement
	measurements = make(map[string][]Measurement)
	var date time.Time
	date = time.Date(2020, 10, 20, 10, 0, 0, 0, time.Local)
	date2 := time.Date(2020, 10, 20, 13, 0, 0, 0, time.Local)
	tagUUIDs := []string{
		"29c3b7d7-d750-4c1a-8695-6001bcff8199",
		"52d2e363-b259-4179-afce-53f463cf58ca",
		"385f326a-8d50-4d5c-a925-010660592ffa",
		"6f351a11-76b9-4cbf-ba1f-ec9a878a59fe",
		"3e032f06-f50a-487a-ae0a-f4768a757581",
		"79815168-f4b9-4544-a2dd-978893072ed9",
		"cce90b30-2607-4d04-99dc-4d0fb1391d50",
		"8dc3480c-36a7-42e2-9d29-d518d435b6ab",
		"35865885-3335-4994-b3c2-c83b2ed64c1e",
		"90d72f4d-9947-4a67-b32c-a5415e94537c",
		"abbbea2a-a36b-4ea8-846e-753a38599ff9",
		"38a401ff-fe0a-487d-9929-34dd4bcfa845",
		"b3dab030-33ba-4fbc-879d-c9b486cb9bd2",
		"35b204b3-67c9-45a1-a3af-402bc4fb3148",
		"eb2c6d4e-06b5-431c-8aa5-bdca136e5918",
		"f031dca9-3014-43f7-b68c-5ed1043e5420",
		"89ca038e-b568-4d2c-ac8d-b2ea0b42e333",
		"850aec8f-5af4-46b0-955c-83fefb5cc962",
		"fb3c2037-f39d-45f2-8d24-1b5e576addcf",
		"27263c4c-e4ce-421c-a152-33eb0d424526",
		"2ca28760-a1ba-43e4-b060-c76b0552efe2",
		"e96d4dda-fc09-43ec-bd0f-6d8781be5f7c",
		"78e421cd-6812-4b50-ab2e-fa1229ba1009",
		"aafa6666-84fb-4d84-bd1a-eea1d2c9164f",
		"62a990d5-2841-47e8-b9d5-8a5c1d9d92c7",
	}
	db = datasource.InitDatabaseConnection("", "5432", "an", "", "")

	for _, uuid := range tagUUIDs {
		var tag Tag
		db.Table("tag").Where("uuid = ?", uuid).First(&tag)
		var measurementsTemp []Measurement
		db.Table("measurement").Where("tag_uuid = ? and date BETWEEN ? AND ?", uuid, date, date2).Find(&measurementsTemp)
		if len(measurementsTemp) > 0 {
			for i, _ := range measurementsTemp {
				measurementsTemp[i].Date = measurementsTemp[i].Date.Add(-3 * time.Hour)
			}
		}
		measurements[tag.Name] = measurementsTemp
	}
	f, err := os.Create("test.csv")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	_, _ = f.WriteString("Имя тега;Дата;Температура;Влажность;\n")
	for k, v := range measurements {
		for _, data := range v {
			_, _ = f.WriteString(fmt.Sprintf("%s;%s;%f;%f;\n", k, data.Date.Format("02.01.2006 15:04:05-0700"), data.Temperature, data.Humidity))
		}
	}
}
