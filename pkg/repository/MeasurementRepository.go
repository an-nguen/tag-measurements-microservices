package repository

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

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
	epsilon float64,
	dataType string) ([]models.Measurement, error) {

	var res []models.Measurement
	if len(uuidList) == 0 || startDate.IsZero() || endDate.IsZero() {
		return res, errors.New("params cannot be nil")
	}

	if err := repo.DataSource.
		Where("tag_uuid::varchar IN ? AND date BETWEEN ? AND ?", uuidList, startDate, endDate).
		Find(&res).Error; err != nil {
		log.Error(err)
	}
	if epsilon > 0.0 {
		if dataType == "" {
			res = utils.DouglasPeuckerMeasurement(res, epsilon, "temperature")
		} else {
			res = utils.DouglasPeuckerMeasurement(res, epsilon, dataType)
		}
	}

	return res, nil
}
