package repository

import (
	"stationhub-api/domain"

	"gorm.io/gorm"
)

type StationRepository struct {
	db *gorm.DB
}

func NewStationRepository(db *gorm.DB) *StationRepository {
	return &StationRepository{db: db}
}

func (r *StationRepository) Create(station *domain.Station) error {
	return r.db.Create(station).Error
}

func (r *StationRepository) Update(station *domain.Station) error {
	return r.db.Save(station).Error
}

func (r *StationRepository) Delete(station *domain.Station) error {
	return r.db.Delete(station).Error
}

func (r *StationRepository) FindByExternalID(externalID string) *domain.Station {
	var station domain.Station
	err := r.db.Where("external_id = ?", externalID).First(&station).Error
	if err != nil {
		return nil
	}
	return &station
}
