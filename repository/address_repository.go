package repository

import (
	"stationhub-api/domain"

	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) Create(address *domain.Address, tx *gorm.DB) error {
	return tx.Create(address).Error
}

func (r *AddressRepository) Update(address *domain.Address, tx *gorm.DB) error {
	return tx.Save(address).Error
}
