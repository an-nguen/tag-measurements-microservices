package repository

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/lib/pq"

	"Thermo-WH/pkg/models"
	"Thermo-WH/pkg/utils"
)

type VoltageTagDataRepository struct {
	DataSource *gorm.DB
}

func (repo VoltageTagDataRepository) GetVoltageTagDataByUUIDs(
	uuidList []string,
	startDate time.Time,
	endDate time.Time,
	epsilon float64) ([]models.Measurement, error) {
	var res []models.Measurement
	if len(uuidList) == 0 || startDate.IsZero() || endDate.IsZero() {
		return res, errors.New("params cannot be nil")
	}

	rows, err := repo.DataSource.DB().Query("SELECT date, voltage, tag_uuid FROM tag_data WHERE tag_uuid = ANY($1) AND voltage IS NOT NULL AND date BETWEEN $2 AND $3 ",
		pq.StringArray(uuidList), startDate, endDate)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var d models.Measurement
		if err := rows.Scan(&d.Date, &d.Voltage, &d.TagUUID); err != nil {
			continue
		}
		res = append(res, d)
	}
	if epsilon > 0.0 {
		res = utils.DouglasPeuckerMeasurement(res, epsilon, "voltage")
	}

	return res, nil
}
