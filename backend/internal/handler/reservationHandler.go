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

type reservationResponse struct {
	ReservationID int64                      `json:"reservationId"`
	Status        string                     `json:"status"`
	Service       reservationServiceResponse `json:"service"`
	SlotID        int64                      `json:"slotId"`
	StartTime     string                     `json:"startTime"`
	EndTime       string                     `json:"endTime"`
	Note          *string                    `json:"note"`
	ReservedAt    string                     `json:"reservedAt"`
	CancelledAt   *string                    `json:"cancelledAt"`
	CreatedAt     string                     `json:"createdAt"`
	UpdatedAt     string                     `json:"updatedAt"`
}

type reservationServiceResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	userID, ok := h.authenticatedUserID(c)
	if !ok {
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

func (h *ReservationHandler) GetReservations(c *gin.Context) {
	userID, ok := h.authenticatedUserID(c)
	if !ok {
		return
	}

	reservations, err := h.reservationService.GetReservationsByUser(c.Request.Context(), userID)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "failed to get reservations")
		return
	}

	response := make([]reservationResponse, 0, len(reservations))
	for i := range reservations {
		response = append(response, newReservationResponse(&reservations[i]))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"reservations": response,
		},
	})
}

func (h *ReservationHandler) GetReservation(c *gin.Context) {
	userID, ok := h.authenticatedUserID(c)
	if !ok {
		return
	}

	reservationID, err := strconv.ParseInt(c.Param("reservationId"), 10, 64)
	if err != nil || reservationID <= 0 {
		h.errorResponse(c, http.StatusBadRequest, "reservationId must be a positive integer")
		return
	}

	reservation, err := h.reservationService.GetReservationDetail(c.Request.Context(), userID, reservationID)
	if err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			h.errorResponse(c, http.StatusNotFound, "reservation not found")
			return
		}

		h.errorResponse(c, http.StatusInternalServerError, "failed to get reservation")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": newReservationResponse(reservation),
	})
}

func (h *ReservationHandler) authenticatedUserID(c *gin.Context) (int64, bool) {
	authClaims, ok := c.Get(middleware.AuthClaimsContextKey)
	if !ok {
		h.errorResponse(c, http.StatusUnauthorized, "authentication context not found")
		return 0, false
	}

	claims, ok := authClaims.(*service.AuthClaims)
	if !ok {
		h.errorResponse(c, http.StatusUnauthorized, "invalid authentication context")
		return 0, false
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		h.errorResponse(c, http.StatusUnauthorized, "invalid user id")
		return 0, false
	}

	return userID, true
}

func newReservationResponse(reservation *service.ReservationResult) reservationResponse {
	var cancelledAt *string
	if reservation.CancelledAt != nil {
		value := reservation.CancelledAt.Format(timeLayout)
		cancelledAt = &value
	}

	return reservationResponse{
		ReservationID: reservation.ReservationID,
		Status:        string(reservation.Status),
		Service: reservationServiceResponse{
			ID:   reservation.ServiceID,
			Name: reservation.ServiceName,
		},
		SlotID:      reservation.SlotID,
		StartTime:   reservation.StartTime.Format(timeLayout),
		EndTime:     reservation.EndTime.Format(timeLayout),
		Note:        reservation.Note,
		ReservedAt:  reservation.ReservedAt.Format(timeLayout),
		CancelledAt: cancelledAt,
		CreatedAt:   reservation.CreatedAt.Format(timeLayout),
		UpdatedAt:   reservation.UpdatedAt.Format(timeLayout),
	}
}

func (h *ReservationHandler) errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}
