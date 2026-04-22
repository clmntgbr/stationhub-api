package domain

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StreetLine1    string    `gorm:"null" json:"street_line_1"`
	StreetLine2    string    `gorm:"null" json:"street_line_2"`
	StreetLine3    string    `gorm:"null" json:"street_line_3"`
	City           string    `gorm:"null" json:"city"`
	State          string    `gorm:"null" json:"state"`
	Zip            string    `gorm:"null" json:"zip"`
	Country        string    `gorm:"null" json:"country"`
	Latitude       float64   `gorm:"null" json:"latitude"`
	Longitude      float64   `gorm:"null" json:"longitude"`
	AdditionalInfo string    `gorm:"null" json:"additional_info"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Address) TableName() string {
	return "addresses"
}
