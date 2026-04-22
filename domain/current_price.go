package domain

import "time"

type CurrentPrice struct {
	Price

	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func (CurrentPrice) TableName() string {
	return "current_prices"
}
