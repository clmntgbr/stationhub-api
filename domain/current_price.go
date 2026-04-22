package domain

import "time"

type CurrentPrice struct {
	Price

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CurrentPrice) TableName() string {
	return "current_prices"
}
