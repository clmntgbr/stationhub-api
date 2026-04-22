package dto

import (
	"stationhub-api/domain"
	"time"
)

type CurrentPriceOutput struct {
	ID        string    `json:"id"`
	Value     float64   `json:"value"`
	Type      string    `json:"type"`
	Currency  string    `json:"currency"`
	Date      time.Time `json:"date"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewCurrentPricesOutput(currentPrices []domain.CurrentPrice) []CurrentPriceOutput {
	outputs := make([]CurrentPriceOutput, len(currentPrices))
	for i, currentPrice := range currentPrices {
		outputs[i] = NewCurrentPriceOutput(currentPrice)
	}
	return outputs
}

func NewCurrentPriceOutput(currentPrice domain.CurrentPrice) CurrentPriceOutput {
	return CurrentPriceOutput{
		ID:        currentPrice.ID.String(),
		Value:     currentPrice.Value,
		Type:      currentPrice.Type,
		Currency:  currentPrice.Currency,
		Date:      currentPrice.Date,
		UpdatedAt: currentPrice.UpdatedAt,
	}
}
