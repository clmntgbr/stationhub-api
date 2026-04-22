package dto

import (
	"stationhub-api/domain"
	"time"
)

type StationOutput struct {
	ID            string               `json:"id"`
	ExternalID    string               `json:"externalId"`
	Name          string               `json:"name"`
	Type          string               `json:"type"`
	Services      []string             `json:"services"`
	Address       AddressOutput        `json:"address"`
	CurrentPrices []CurrentPriceOutput `json:"currentPrices"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}

func NewStationOutput(station domain.Station) StationOutput {
	return StationOutput{
		ID:            station.ID.String(),
		ExternalID:    station.ExternalID,
		Name:          station.Name,
		Type:          station.Type,
		Services:      station.Services,
		Address:       NewAddressOutput(station.Address),
		CurrentPrices: NewCurrentPricesOutput(station.CurrentPrices),
		CreatedAt:     station.CreatedAt,
		UpdatedAt:     station.UpdatedAt,
	}
}

func NewMinimalStationsOutput(stations []domain.Station) []StationOutput {
	outputs := make([]StationOutput, len(stations))
	for i, station := range stations {
		outputs[i] = NewStationOutput(station)
	}
	return outputs
}
