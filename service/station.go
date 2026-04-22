package service

import (
	"stationhub-api/dto"
	"stationhub-api/repository"
)

type StationService struct {
	stationRepository *repository.StationRepository
}

func NewStationService(stationRepository *repository.StationRepository) *StationService {
	return &StationService{
		stationRepository: stationRepository,
	}
}

func (s *StationService) GetStations() ([]dto.StationOutput, error) {
	stations, err := s.stationRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.NewMinimalStationsOutput(stations), nil
}
