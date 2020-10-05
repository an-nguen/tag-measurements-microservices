package utils

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"Thermo-WH/pkg/datasource"
	"Thermo-WH/pkg/repository"
)

var db = datasource.InitDatabaseConnection("127.0.0.1", "5432", "an",
	"FoX13@", "temperature_display_db")

var tagTemperatureRepo = repository.MeasurementRepository{
	DataSource: db,
}

func TestDouglasPeucker(t *testing.T) {
	var uuidList []string
	uuidList = append(uuidList, "3e032f06-f50a-487a-ae0a-f4768a757581")
	res, _ := tagTemperatureRepo.GetMeasurementByUUIDs(uuidList, time.Now().Add(time.Hour*24*180), time.Now(), 0)
	fmt.Printf("len %d\n", len(res))
	resPnt := DouglasPeuckerMeasurement(res, 0, "temperature")
	fmt.Printf("len %d\n", len(resPnt))
	if len(res) < len(resPnt) {
		t.Error("Expected that length of res less than length of resPnt, got difference ", len(resPnt)-len(res))
	}
}
