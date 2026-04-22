package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"stationhub-api/domain"
	"stationhub-api/dto"
	"stationhub-api/repository"
	"strings"
	"sync"

	"golang.org/x/text/encoding/charmap"
)

const (
	WorkerCount = 8
	JobQueueSize = 100
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

	address := &domain.Address{
		StreetLine1: pdv.Adresse,
		City:        pdv.Ville,
		State:       pdv.CP,
		Country:     "France",
		Latitude:    pdv.Latitude,
		Longitude:   pdv.Longitude,
	}

	station := &domain.Station{
		ExternalID: pdv.ID,
		Name:       pdv.Adresse + " " + pdv.Ville,
		Type:       "gas",
	}

	created, err := s.stationRepository.CreateStationWithAddress(station, address, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create station with address: %w", err)
	}

	if !created {
		tx.Rollback()
		return nil
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
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
