package dto

import (
	"stationhub-api/domain"
	"time"
)

type CurrentPriceOutput struct {
	MinimalCurrentPriceOutput MinimalCurrentPriceOutput `json:"minimalCurrentPriceOutput"`
	ID                        string                    `json:"id"`
	Currency                  string                    `json:"currency"`
	UpdatedAt                 time.Time                 `json:"updatedAt"`
}

type MinimalCurrentPriceOutput struct {
	Value float64   `json:"value"`
	Type  string    `json:"type"`
	Date  time.Time `json:"date"`
}

func NewCurrentPricesOutput(currentPrices []domain.CurrentPrice) []CurrentPriceOutput {
	outputs := make([]CurrentPriceOutput, len(currentPrices))
	for i, currentPrice := range currentPrices {
		outputs[i] = NewCurrentPriceOutput(currentPrice)
	}
	return outputs
}

func NewMinimalCurrentPricesOutput(currentPrices []domain.CurrentPrice) []MinimalCurrentPriceOutput {
	outputs := make([]MinimalCurrentPriceOutput, len(currentPrices))
	for i, currentPrice := range currentPrices {
		outputs[i] = NewMinimalCurrentPriceOutput(currentPrice)
	}
	return outputs
}

func NewMinimalCurrentPriceOutput(currentPrice domain.CurrentPrice) MinimalCurrentPriceOutput {
	return MinimalCurrentPriceOutput{
		Value: currentPrice.Value,
		Type:  currentPrice.Type,
		Date:  currentPrice.Date,
	}
}

func NewCurrentPriceOutput(currentPrice domain.CurrentPrice) CurrentPriceOutput {
	return CurrentPriceOutput{
		MinimalCurrentPriceOutput: NewMinimalCurrentPriceOutput(currentPrice),
		ID:                        currentPrice.ID.String(),
		Currency:                  currentPrice.Currency,
		UpdatedAt:                 currentPrice.UpdatedAt,
	}
}
