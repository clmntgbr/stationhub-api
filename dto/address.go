package dto

import (
	"stationhub-api/domain"
)

type AddressOutput struct {
	MinimalAddressOutput MinimalAddressOutput `json:"minimalAddressOutput"`
	ID                   string               `json:"id"`
	StreetLine1          string               `json:"streetLine1"`
	StreetLine2          string               `json:"streetLine2"`
	StreetLine3          string               `json:"streetLine3"`
	City                 string               `json:"city"`
	State                string               `json:"state"`
	Zip                  string               `json:"zip"`
	Country              string               `json:"country"`
	AdditionalInfo       string               `json:"additionalInfo"`
}

type MinimalAddressOutput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewMinimalAddressOutput(address domain.Address) MinimalAddressOutput {
	return MinimalAddressOutput{
		Latitude:  address.Latitude,
		Longitude: address.Longitude,
	}
}

func NewAddressOutput(address domain.Address) AddressOutput {
	return AddressOutput{
		MinimalAddressOutput: NewMinimalAddressOutput(address),
		ID:                   address.ID.String(),
		StreetLine1:          address.StreetLine1,
		StreetLine2:          address.StreetLine2,
		StreetLine3:          address.StreetLine3,
		City:                 address.City,
		State:                address.State,
		Zip:                  address.Zip,
		Country:              address.Country,
		AdditionalInfo:       address.AdditionalInfo,
	}
}
