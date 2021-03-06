package repository

import (
	"gorm.io/gorm"

	"tag-measurements-microservices/pkg/models"
)

type TagRepository struct {
	DataSource *gorm.DB
}

func (tr TagRepository) GetTags() ([]models.Tag, error) {
	var tags []models.Tag
	if err := tr.DataSource.Find(&tags).Error; err != nil {
		return tags, err
	}
	return tags, nil
}

func (tr TagRepository) GetTagsByMAC(mac string) ([]models.Tag, error) {
	var tags []models.Tag
	if err := tr.DataSource.Find(&tags, "mac_tag_manager = ?", mac).Error; err != nil {
		return tags, err
	}
	return tags, nil
}

func (tr TagRepository) GetTagsByTemperatureZone(id string) ([]models.Tag, error) {
	var tags []models.Tag
	var temperatureZone models.TemperatureZone
	if err := tr.DataSource.Preload("Tags").First(&temperatureZone, id).Error; err != nil {
		return nil, err
	}
	tr.DataSource.Preload("Tags").Model(&temperatureZone)
	for _, tag := range temperatureZone.Tags {
		tags = append(tags, *tag)
	}

	return tags, nil
}

func (tr TagRepository) GetTagsByUUID(uuid string) ([]models.Tag, error) {
	var tags []models.Tag
	if err := tr.DataSource.Find(tags, "uuid = ?", uuid).Error; err != nil {
		return tags, err
	}

	return tags, nil
}

func (tr TagRepository) UpdateTag(tag models.Tag, uuid string) (models.Tag, error) {
	var tagDb models.Tag
	tr.DataSource.Preload("TemperatureZones").Where("uuid = ?", uuid).First(&tagDb)
	tagDb.VerificationDate = tag.VerificationDate
	tr.DataSource.Updates(&tagDb)

	return tagDb, nil
}
