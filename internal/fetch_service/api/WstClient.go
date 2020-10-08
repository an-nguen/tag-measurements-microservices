package api

import (
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
}
