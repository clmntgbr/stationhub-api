package repository

import (
	"stationhub-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *StationRepository) CreateStationWithAddress(station *domain.Station, address *domain.Address, tx *gorm.DB) (uuid.UUID, error) {
	var existingStation domain.Station
	err := tx.Select("id", "external_id", "address_id").Where("external_id = ?", station.ExternalID).First(&existingStation).Error

	if err == nil {
		station.ID = existingStation.ID
		station.AddressID = existingStation.AddressID
		
		result := tx.Omit("Address").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"services", "name", "type", "updated_at"}),
		}).Create(station)

		return existingStation.ID, result.Error
	}

	if err := tx.Create(address).Error; err != nil {
		return uuid.Nil, err
	}

	station.AddressID = address.ID

	result := tx.Omit("Address").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"services", "name", "type", "updated_at"}),
	}).Create(station)

	return station.ID, result.Error
}

func (r *StationRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}
