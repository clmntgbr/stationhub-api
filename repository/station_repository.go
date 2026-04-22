package repository

import (
	"stationhub-api/domain"

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

func (r *StationRepository) CreateStationWithAddress(station *domain.Station, address *domain.Address, tx *gorm.DB) (bool, error) {
	if err := tx.Create(address).Error; err != nil {
		return false, err
	}

	station.AddressID = address.ID

	result := tx.Omit("Address").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoNothing: true,
	}).Create(station)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (r *StationRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}
