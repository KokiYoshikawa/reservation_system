package handler

import (
	"errors"
	"net/http"
	"strconv"

	"reservation-system/backend/internal/dto"
	"reservation-system/backend/internal/middleware"
	"reservation-system/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	reservationService service.ReservationService
}

type createReservationRequest struct {
	ServiceID int64  `json:"serviceId"`
	SlotID    int64  `json:"slotId"`
	Note      string `json:"note"`
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	authClaims, ok := c.Get(middleware.AuthClaimsContextKey)
	if !ok {
		h.errorResponse(c, http.StatusUnauthorized, "authentication context not found")
		return
	}

	claims, ok := authClaims.(*service.AuthClaims)
	if !ok {
		h.errorResponse(c, http.StatusUnauthorized, "invalid authentication context")
		return
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		h.errorResponse(c, http.StatusUnauthorized, "invalid user id")
		return
	}

	var req createReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ServiceID <= 0 || req.SlotID <= 0 {
		h.errorResponse(c, http.StatusBadRequest, "serviceId and slotId must be positive integers")
		return
	}

	result, err := h.reservationService.CreateReservation(c.Request.Context(), userID, dto.CreateReservationRequest{
		ServiceID: req.ServiceID,
		SlotID:    req.SlotID,
		Note:      req.Note,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrServiceNotFound):
			h.errorResponse(c, http.StatusNotFound, "service not found")
			return
		case errors.Is(err, service.ErrReservationSlotNotFound):
			h.errorResponse(c, http.StatusNotFound, "reservation slot not found")
			return
		case errors.Is(err, service.ErrReservationSlotPast):
			h.errorResponse(c, http.StatusBadRequest, "reservation slot is in the past")
			return
		case errors.Is(err, service.ErrReservationConflict):
			h.errorResponse(c, http.StatusConflict, "reservation slot is already reserved")
			return
		default:
			h.errorResponse(c, http.StatusInternalServerError, "failed to create reservation")
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"reservationId": result.ReservationID,
			"status":        result.Status,
			"service": gin.H{
				"id":   result.ServiceID,
				"name": result.ServiceName,
			},
			"startTime":  result.StartTime.Format(timeLayout),
			"endTime":    result.EndTime.Format(timeLayout),
			"reservedAt": result.ReservedAt.Format(timeLayout),
		},
	})
}

func (h *ReservationHandler) errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}
