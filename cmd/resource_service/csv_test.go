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
	date = time.Now().Add(-7 * (time.Hour * 24))
	tagUUIDs := []string{"f031dca9-3014-43f7-b68c-5ed1043e5420", "eb2c6d4e-06b5-431c-8aa5-bdca136e5918", "62a990d5-2841-47e8-b9d5-8a5c1d9d92c7",
		"aafa6666-84fb-4d84-bd1a-eea1d2c9164f", "850aec8f-5af4-46b0-955c-83fefb5cc962", "89ca038e-b568-4d2c-ac8d-b2ea0b42e333", "27263c4c-e4ce-421c-a152-33eb0d424526",
		"2ca28760-a1ba-43e4-b060-c76b0552efe2", "78e421cd-6812-4b50-ab2e-fa1229ba1009", "e96d4dda-fc09-43ec-bd0f-6d8781be5f7c", "fb3c2037-f39d-45f2-8d24-1b5e576addcf"}
	db = datasource.InitDatabaseConnection("", "", "", "", "")

	for _, uuid := range tagUUIDs {
		var tag Tag
		db.Table("tag").Where("uuid = ?", uuid).First(&tag)
		var measurementsTemp []Measurement
		db.Table("tag_data").Where("tag_uuid = ? and date BETWEEN ? AND ?", uuid, date, time.Now()).Find(&measurementsTemp)
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
			_, _ = f.WriteString(fmt.Sprintf("%s;%s;%f;%f;\n", k, data.Date.Format("01.02.2006 15:04:05-0700"), data.Temperature, data.Humidity))
		}
	}
}
