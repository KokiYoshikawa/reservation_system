package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"reservation-system/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ReservationSlotHandler struct {
	reservationSlotService service.ReservationSlotService
}

type reservationSlotResponse struct {
	SlotID        int64  `json:"slotId"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	Capacity      int    `json:"capacity"`
	ReservedCount int    `json:"reservedCount"`
	Available     bool   `json:"available"`
}

func NewReservationSlotHandler(reservationSlotService service.ReservationSlotService) *ReservationSlotHandler {
	return &ReservationSlotHandler{reservationSlotService: reservationSlotService}
}

func (h *ReservationSlotHandler) SearchAvailableSlots(c *gin.Context) {
	dateValue := c.Query("date")
	serviceIDValue := c.Query("serviceId")
	if dateValue == "" || serviceIDValue == "" {
		h.errorResponse(c, http.StatusBadRequest, "date and serviceId are required")
		return
	}

	targetDate, err := time.Parse("2006-01-02", dateValue)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "date must be in YYYY-MM-DD format")
		return
	}

	serviceID, err := strconv.ParseInt(serviceIDValue, 10, 64)
	if err != nil || serviceID <= 0 {
		h.errorResponse(c, http.StatusBadRequest, "serviceId must be a positive integer")
		return
	}

	result, err := h.reservationSlotService.SearchAvailableSlots(c.Request.Context(), targetDate, serviceID)
	if err != nil {
		if errors.Is(err, service.ErrServiceNotFound) {
			h.errorResponse(c, http.StatusNotFound, "service not found")
			return
		}

		h.errorResponse(c, http.StatusInternalServerError, "failed to search reservation slots")
		return
	}

	slots := make([]reservationSlotResponse, 0, len(result.Slots))
	for _, slot := range result.Slots {
		slots = append(slots, reservationSlotResponse{
			SlotID:        slot.SlotID,
			StartTime:     slot.StartTime.Format(timeLayout),
			EndTime:       slot.EndTime.Format(timeLayout),
			Capacity:      slot.Capacity,
			ReservedCount: slot.ReservedCount,
			Available:     slot.Available,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"date":  result.Date,
			"slots": slots,
		},
	})
}

func (h *ReservationSlotHandler) errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}
