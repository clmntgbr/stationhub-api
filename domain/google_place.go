package domain

import (
	"time"

	"github.com/google/uuid"
)

type GooglePlace struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	PlaceID                  string  `gorm:"null" json:"place_id"`
	InternationalPhoneNumber string  `gorm:"null" json:"international_phone_number"`
	Rating                   float64 `gorm:"null" json:"rating"`
	UserRatingCount          int     `gorm:"null" json:"user_rating_count"`
	BusinessStatus           string  `gorm:"null" json:"business_status"`
	WebsiteURL               string  `gorm:"null" json:"website_url"`
	GoogleMapDirectionsURL   string  `gorm:"null" json:"google_map_directions_url"`
	GoogleMapURL             string  `gorm:"null" json:"google_map_url"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (GooglePlace) TableName() string {
	return "google_places"
}
