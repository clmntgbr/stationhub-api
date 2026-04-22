package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"stationhub-api/domain"
	"stationhub-api/dto"
	"stationhub-api/repository"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type GasPricesUpdateService struct {
	stationRepository *repository.StationRepository
}

func NewGasPricesUpdateService(stationRepository *repository.StationRepository) *GasPricesUpdateService {
	return &GasPricesUpdateService{stationRepository: stationRepository}
}

func (s *GasPricesUpdateService) UpdateGasPrices(xmlFilePath string) error {
	pdvListe, err := s.ExtractGasPricesFile(xmlFilePath)
	if err != nil {
		return fmt.Errorf("failed to extract gas prices file: %w", err)
	}

	for i := 0; i < len(pdvListe.PDVs); i++ {
		station := s.stationRepository.FindByExternalID(pdvListe.PDVs[i].ID)
		if station == nil {
			station = &domain.Station{
				ExternalID: pdvListe.PDVs[i].ID,
				Name:       pdvListe.PDVs[i].Adresse + " " + pdvListe.PDVs[i].Ville,
				Type:       "gas",
				Address: domain.Address{
					StreetLine1: pdvListe.PDVs[i].Adresse,
					City:        pdvListe.PDVs[i].Ville,
					State:       pdvListe.PDVs[i].CP,
					Country:     "France",
					Latitude:    pdvListe.PDVs[i].Latitude,
					Longitude:   pdvListe.PDVs[i].Longitude,
				},
			}
			s.stationRepository.Create(station)
		}
	}

	return nil
}

func (s *GasPricesUpdateService) ExtractGasPricesFile(xmlFilePath string) (dto.PDVListe, error) {
	xmlFile, err := os.Open(xmlFilePath)
	if err != nil {
		return dto.PDVListe{}, fmt.Errorf("failed to open XML file: %w", err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch strings.ToLower(charset) {
		case "iso-8859-1", "latin1":
			return charmap.ISO8859_1.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unsupported charset: %s", charset)
		}
	}

	var pdvListe dto.PDVListe
	err = decoder.Decode(&pdvListe)
	if err != nil {
		return dto.PDVListe{}, fmt.Errorf("failed to unmarshal XML file: %w", err)
	}

	return pdvListe, nil
}
