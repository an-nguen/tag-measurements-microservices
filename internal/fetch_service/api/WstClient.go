package api

import (
	"github.com/jinzhu/gorm"
	"net/http"

	"tag-measurements-microservices/pkg/models"
)

type WstClient struct {
	Client               *http.Client
	MeasurementTableName string
	HostUrl              string
	SessionId            string
	TagManager           models.TagManager
	Email                string
	Password             string
	Connection           *gorm.DB
}
