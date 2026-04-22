package handler

import (
	"stationhub-api/service"

	"github.com/gofiber/fiber/v3"
)

type StationHandler struct {
	BaseHandler
	stationService *service.StationService
}

func NewStationHandler(stationService *service.StationService) *StationHandler {
	return &StationHandler{
		stationService: stationService,
	}
}

func (h *StationHandler) GetStations(c fiber.Ctx) error {
	stations, err := h.stationService.GetStations()
	if err != nil {
		return h.sendInternalError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(stations)
}
