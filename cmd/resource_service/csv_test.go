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
	date = time.Now().Add(-14 * (time.Hour * 24))
	tagUUIDs := []string{"92bb22cf-7273-4b98-833f-dc90782dee42", "61ac408f-6e5c-4373-b8fd-919563f62352",
		"0b91bf39-d8c3-44c4-83c1-5162dc156483",
		"f8a5de85-a6f5-4c38-87ec-af8d113a1b6e",
		"d1fd7a4f-26ce-4026-aef0-3f2aa65b284b",
		"a0da30ae-0db6-404a-8dff-476176eec409",
		"6e1e132e-135b-4a79-b7e9-45ab4c2a7641",
		"ee0690f8-420f-4662-b85d-2c064c72c996",
		"92bb22cf-7273-4b98-833f-dc90782dee42",
		"62ed2142-3bf1-4f1f-a517-d27b11e22cb2",
	}
	db = datasource.InitDatabaseConnection("", "5432", "an", "", "tag_measurements")

	for _, uuid := range tagUUIDs {
		var tag Tag
		db.Table("tag").Where("uuid = ?", uuid).First(&tag)
		var measurementsTemp []Measurement
		db.Table("measurement").Where("tag_uuid = ? and date BETWEEN ? AND ?", uuid, date, time.Now()).Find(&measurementsTemp)
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
