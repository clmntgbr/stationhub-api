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

func (r *StationRepository) FindByExternalIDWithTx(externalID string, tx *gorm.DB) (*domain.Station, error) {
	var station domain.Station
	err := tx.Select("id", "external_id", "address_id").Where("external_id = ?", externalID).First(&station).Error
	if err != nil {
		return nil, err
	}
	return &station, nil
}

func (r *StationRepository) CreateWithTx(station *domain.Station, tx *gorm.DB) error {
	return tx.Omit("Address").Create(station).Error
}

func (r *StationRepository) UpdateWithTx(station *domain.Station, tx *gorm.DB) error {
	return tx.Omit("Address").Updates(station).Error
}

func (r *StationRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *StationRepository) FindAll() ([]domain.Station, error) {
	var stations []domain.Station
	err := r.db.Preload("Address").Preload("CurrentPrices").Limit(10).Find(&stations).Error
	if err != nil {
		return nil, err
	}
	return stations, nil
}
