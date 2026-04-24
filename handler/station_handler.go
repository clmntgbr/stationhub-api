package handler

import (
	"stationhub-api/dto"
	"stationhub-api/errors"
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
	var query dto.GetStationsQuery

	if err := c.Bind().Query(&query); err != nil {
		return h.sendBadRequest(c, errors.ErrInvalidRequestBody)
	}

	stations, err := h.stationService.GetStations(query)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(stations)
}
