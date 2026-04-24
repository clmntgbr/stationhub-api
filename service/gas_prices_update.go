package service

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"stationhub-api/domain"
	"stationhub-api/dto"
	"stationhub-api/repository"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/encoding/charmap"
	"gorm.io/gorm"
)

const (
	WorkerCount  = 8
	JobQueueSize = 100
)

type GasPricesUpdateService struct {
	stationRepository      *repository.StationRepository
	addressRepository      *repository.AddressRepository
	currentPriceRepository *repository.CurrentPriceRepository
	priceHistoryRepository *repository.PriceHistoryRepository
}

func NewGasPricesUpdateService(
	stationRepository *repository.StationRepository,
	addressRepository *repository.AddressRepository,
	currentPriceRepository *repository.CurrentPriceRepository,
	priceHistoryRepository *repository.PriceHistoryRepository,
) *GasPricesUpdateService {
	return &GasPricesUpdateService{
		stationRepository:      stationRepository,
		addressRepository:      addressRepository,
		currentPriceRepository: currentPriceRepository,
		priceHistoryRepository: priceHistoryRepository,
	}
}

func (s *GasPricesUpdateService) UpdateGasPrices(xmlFilePath string) error {
	pdvListe, err := s.ExtractGasPricesFile(xmlFilePath)
	if err != nil {
		return fmt.Errorf("failed to extract gas prices file: %w", err)
	}

	totalPDVs := len(pdvListe.PDVs)
	log.Printf("Starting ingestion of %d stations with %d workers", totalPDVs, WorkerCount)

	jobQueue := make(chan dto.PDV, JobQueueSize)
	var wg sync.WaitGroup

	for i := 0; i < WorkerCount; i++ {
		wg.Add(1)
		go s.worker(i+1, jobQueue, &wg)
	}

	for _, pdv := range pdvListe.PDVs {
		jobQueue <- pdv
	}
	close(jobQueue)

	wg.Wait()

	log.Printf("Ingestion completed for %d stations", totalPDVs)
	return nil
}

func (s *GasPricesUpdateService) worker(id int, jobs <-chan dto.PDV, wg *sync.WaitGroup) {
	defer wg.Done()

	processed := 0
	for pdv := range jobs {
		if err := s.processPDV(pdv); err != nil {
			log.Printf("[Worker %d] Error processing station %s: %v", id, pdv.ID, err)
		}
		processed++
	}

	log.Printf("[Worker %d] Processed %d stations", id, processed)
}

func (s *GasPricesUpdateService) processPDV(pdv dto.PDV) error {
	tx := s.stationRepository.BeginTransaction()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Panic recovered for station %s: %v", pdv.ID, r)
		}
	}()

	existingStation, err := s.stationRepository.FindByExternalIDWithTx(pdv.ID, tx)

	var stationID uuid.UUID

	if err == nil {
		station := &domain.Station{
			ID:         existingStation.ID,
			ExternalID: pdv.ID,
			Name:       pdv.Adresse + " " + pdv.Ville,
			Type:       "gas",
			Services:   pdv.Services.List,
			AddressID:  existingStation.AddressID,
		}

		if err := s.stationRepository.UpdateWithTx(station, tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update station: %w", err)
		}

		stationID = existingStation.ID
	} else {
		address := &domain.Address{
			StreetLine1: pdv.Adresse,
			City:        pdv.Ville,
			State:       pdv.CP,
			Country:     "France",
			Latitude:    pdv.Latitude / 100000,
			Longitude:   pdv.Longitude / 100000,
		}

		if err := s.addressRepository.Create(address, tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create address: %w", err)
		}

		station := &domain.Station{
			ExternalID: pdv.ID,
			Name:       pdv.Adresse + " " + pdv.Ville,
			Type:       "gas",
			Services:   pdv.Services.List,
			AddressID:  address.ID,
		}

		if err := s.stationRepository.CreateWithTx(station, tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create station: %w", err)
		}

		stationID = station.ID
	}

	if err := s.processPrices(tx, stationID, pdv.Prix); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to process prices: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *GasPricesUpdateService) processPrices(tx *gorm.DB, stationID uuid.UUID, prices []dto.Prix) error {
	for _, prix := range prices {
		priceDate, err := time.Parse("2006-01-02 15:04:05", prix.Maj)
		if err != nil {
			return fmt.Errorf("invalid date format for prix %s: %w", prix.Nom, err)
		}

		existingPrice, err := s.currentPriceRepository.FindByStationAndType(stationID, prix.ID, tx)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to find current price: %w", err)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newPrice := &domain.CurrentPrice{
				Price: domain.Price{
					StationID: stationID,
					Type:      prix.Nom,
					TypeId:    prix.ID,
					Value:     prix.Valeur,
					Currency:  "EUR",
					Date:      priceDate,
				},
			}
			if err := s.currentPriceRepository.Create(newPrice, tx); err != nil {
				return fmt.Errorf("failed to create current price: %w", err)
			}

			history := &domain.PriceHistory{
				Price: domain.Price{
					StationID: newPrice.StationID,
					Type:      newPrice.Type,
					TypeId:    newPrice.TypeId,
					Value:     newPrice.Value,
					Currency:  newPrice.Currency,
					Date:      newPrice.Date,
				},
			}
			if err := s.priceHistoryRepository.Create(history, tx); err != nil {
				return fmt.Errorf("failed to create price history: %w", err)
			}

			continue
		}

		if !priceDate.After(existingPrice.Date) {
			continue
		}

		if existingPrice.Value == prix.Valeur {
			existingPrice.Date = priceDate
			if err := s.currentPriceRepository.Update(existingPrice, tx); err != nil {
				return fmt.Errorf("failed to update current price date: %w", err)
			}
			continue
		}

		history := &domain.PriceHistory{
			Price: domain.Price{
				StationID: existingPrice.StationID,
				Type:      existingPrice.Type,
				TypeId:    existingPrice.TypeId,
				Value:     existingPrice.Value,
				Currency:  existingPrice.Currency,
				Date:      existingPrice.Date,
			},
		}

		if err := s.priceHistoryRepository.Create(history, tx); err != nil {
			return fmt.Errorf("failed to create price history: %w", err)
		}

		existingPrice.Value = prix.Valeur
		existingPrice.Date = priceDate

		if err := s.currentPriceRepository.Update(existingPrice, tx); err != nil {
			return fmt.Errorf("failed to update current price: %w", err)
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
