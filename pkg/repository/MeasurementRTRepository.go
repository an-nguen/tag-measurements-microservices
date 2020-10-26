package repository

import (
	"errors"
	"tag-measurements-microservices/pkg/utils"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"tag-measurements-microservices/pkg/models"
)

type MeasurementRTRepository struct {
	DataSource *gorm.DB
}

func (repo MeasurementRTRepository) GetMeasurement() ([]models.MeasurementRT, error) {
	var res []models.MeasurementRT
	if err := repo.DataSource.Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (repo MeasurementRTRepository) GetMeasurementByUUIDs(
	uuidList []string,
	startDate time.Time,
	endDate time.Time,
	epsilon float64,
	dataType string) ([]models.MeasurementRT, error) {

	var res []models.MeasurementRT
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
			res = utils.DouglasPeuckerMeasurementRT(res, epsilon, "temperature")
		} else {
			res = utils.DouglasPeuckerMeasurementRT(res, epsilon, dataType)
		}
	}

	return res, nil
}
