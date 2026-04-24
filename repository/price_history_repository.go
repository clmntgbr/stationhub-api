package repository

import (
	"stationhub-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PriceHistoryRepository struct {
	db *gorm.DB
}

func NewPriceHistoryRepository(db *gorm.DB) *PriceHistoryRepository {
	return &PriceHistoryRepository{db: db}
}

func (r *PriceHistoryRepository) FindByStationAndType(stationID uuid.UUID, priceTypeId int, tx *gorm.DB) (*domain.PriceHistory, error) {
	var price domain.PriceHistory
	err := tx.Where("station_id = ? AND type_id = ?", stationID, priceTypeId).First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *PriceHistoryRepository) Create(price *domain.PriceHistory, tx *gorm.DB) error {
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(price).Error
}

func (r *PriceHistoryRepository) Update(price *domain.PriceHistory, tx *gorm.DB) error {
	return tx.Save(price).Error
}
