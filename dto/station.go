package dto

import (
	"stationhub-api/domain"
	"time"
)

type StationOutput struct {
	MinimalStationOutput MinimalStationOutput `json:"minimalStationOutput"`
	ID                   string               `json:"id"`
	Type                 string               `json:"type"`
	Services             []string             `json:"services"`
	CreatedAt            time.Time            `json:"createdAt"`
	UpdatedAt            time.Time            `json:"updatedAt"`
}

type MinimalStationOutput struct {
	ExternalID    string                      `json:"externalId"`
	Name          string                      `json:"name"`
	Address       MinimalAddressOutput        `json:"address"`
	CurrentPrices []MinimalCurrentPriceOutput `json:"currentPrices"`
}

func NewStationOutput(station domain.Station) StationOutput {
	return StationOutput{
		ID:                   station.ID.String(),
		MinimalStationOutput: NewMinimalStationOutput(station),
		CreatedAt:            station.CreatedAt,
		UpdatedAt:            station.UpdatedAt,
	}
}

func NewMinimalStationOutput(station domain.Station) MinimalStationOutput {
	return MinimalStationOutput{
		ExternalID:    station.ExternalID,
		Name:          station.Name,
		Address:       NewMinimalAddressOutput(station.Address),
		CurrentPrices: NewMinimalCurrentPricesOutput(station.CurrentPrices),
	}
}

func NewMinimalStationsOutput(stations []domain.Station) []MinimalStationOutput {
	outputs := make([]MinimalStationOutput, len(stations))
	for i, station := range stations {
		outputs[i] = NewMinimalStationOutput(station)
	}
	return outputs
}

func NewStationsOutput(stations []domain.Station) []StationOutput {
	outputs := make([]StationOutput, len(stations))
	for i, station := range stations {
		outputs[i] = NewStationOutput(station)
	}
	return outputs
}
