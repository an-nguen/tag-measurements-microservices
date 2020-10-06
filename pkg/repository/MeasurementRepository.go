package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"

	"tag-measurements-microservices/pkg/models"
	"tag-measurements-microservices/pkg/utils"
)

type MeasurementRepository struct {
	DataSource *gorm.DB
}

func (repo MeasurementRepository) GetMeasurement() ([]models.Measurement, error) {
	var res []models.Measurement
	if err := repo.DataSource.Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (repo MeasurementRepository) GetMeasurementByUUIDs(
	uuidList []string,
	startDate time.Time,
	endDate time.Time,
	epsilon float64) ([]models.Measurement, error) {

	var res []models.Measurement
	if len(uuidList) == 0 || startDate.IsZero() || endDate.IsZero() {
		return res, errors.New("params cannot be nil")
	}

	rows, err := repo.DataSource.DB().Query("SELECT date, temperature, humidity, voltage, signal, tag_uuid FROM measurement WHERE tag_uuid = ANY($1) AND date BETWEEN $2 AND $3",
		pq.StringArray(uuidList), startDate, endDate)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var d models.Measurement
		if err := rows.Scan(&d.Date, &d.Temperature, &d.Humidity, &d.Voltage, &d.Signal, &d.TagUUID); err != nil {
			continue
		}
		res = append(res, d)
	}
	if epsilon > 0.0 {
		res = utils.DouglasPeuckerMeasurement(res, epsilon, "temperature")
	}

	return res, nil
}
