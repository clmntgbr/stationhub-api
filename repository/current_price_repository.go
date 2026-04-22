package repository

import (
	"stationhub-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CurrentPriceRepository struct {
	db *gorm.DB
}

func NewCurrentPriceRepository(db *gorm.DB) *CurrentPriceRepository {
	return &CurrentPriceRepository{db: db}
}

func (r *CurrentPriceRepository) FindByStationAndType(stationID uuid.UUID, priceType string, tx *gorm.DB) (*domain.CurrentPrice, error) {
	var price domain.CurrentPrice
	err := tx.Where("station_id = ? AND type = ?", stationID, priceType).First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *CurrentPriceRepository) Create(price *domain.CurrentPrice, tx *gorm.DB) error {
	return tx.Create(price).Error
}

func (r *CurrentPriceRepository) Update(price *domain.CurrentPrice, tx *gorm.DB) error {
	return tx.Save(price).Error
}
