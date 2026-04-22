package domain

import (
	"time"

	"github.com/google/uuid"
)

type Station struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ExternalID string    `gorm:"uniqueIndex;not null" json:"external_id"`
	Name       string    `gorm:"null" json:"name"`
	Type       string    `gorm:"null" json:"type"`

	AddressID uuid.UUID `gorm:"type:uuid;not null" json:"address_id"`
	Address   Address   `gorm:"foreignKey:AddressID" json:"address"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Station) TableName() string {
	return "stations"
}
