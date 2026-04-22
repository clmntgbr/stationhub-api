package domain

import (
	"time"

	"github.com/google/uuid"
)

type CurrentPrice struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Value    float64   `gorm:"not null" json:"value"`
	Currency string    `gorm:"not null" json:"currency"`
	Type     string    `gorm:"not null" json:"type"`
	Date     time.Time `gorm:"not null" json:"date"`

	StationID uuid.UUID `gorm:"type:uuid;not null" json:"station_id"`
	Station   Station   `gorm:"foreignKey:StationID" json:"station"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CurrentPrice) TableName() string {
	return "current_prices"
}
