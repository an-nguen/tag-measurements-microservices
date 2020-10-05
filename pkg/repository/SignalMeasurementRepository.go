package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/lib/pq"

	"Thermo-WH/pkg/models"
	"Thermo-WH/pkg/utils"
)

type SignalTagDataRepository struct {
	DataSource *gorm.DB
}

func (repo SignalTagDataRepository) GetSignalTagDataByUUIDs(
	uuidList []string, startDate time.Time,
	endDate time.Time, epsilon float64) ([]models.Measurement, error) {
	var res []models.Measurement
	if len(uuidList) == 0 || startDate.IsZero() || endDate.IsZero() {
		return res, errors.New("params cannot be nil")
	}

	rows, err := repo.DataSource.DB().Query("SELECT date, signal, tag_uuid FROM tag_data WHERE tag_uuid = ANY($1) AND signal IS NOT NULL AND date BETWEEN $2 AND $3",
		pq.StringArray(uuidList), startDate, endDate)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var d models.Measurement
		if err := rows.Scan(&d.Date, &d.Signal, &d.TagUUID); err != nil {
			continue
		}
		res = append(res, d)
	}
	if epsilon > 0.0 {
		res = utils.DouglasPeuckerMeasurement(res, epsilon, "signal")
	}

	return res, nil
}
