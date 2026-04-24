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

func (s *StationService) GetStations(q dto.GetStationsQuery) ([]dto.MinimalStationOutput, error) {
	stations, err := s.stationRepository.FindNearby(q.Latitude, q.Longitude, q.Radius)
	if err != nil {
		return nil, err
	}
	return dto.NewMinimalStationsOutput(stations), nil
}
