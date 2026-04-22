package domain

import "time"

type PriceHistory struct {
	Price

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PriceHistory) TableName() string {
	return "price_histories"
}
