package handler

import (
	"net/http"

	"reservation-system/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	serviceService service.ServiceService
}

type serviceResponse struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	DurationMinutes int    `json:"duration_minutes"`
	Price           int    `json:"price"`
	IsActive        bool   `json:"is_active"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func NewServiceHandler(serviceService service.ServiceService) *ServiceHandler {
	return &ServiceHandler{serviceService: serviceService}
}

func (h *ServiceHandler) ListActiveServices(c *gin.Context) {
	services, err := h.serviceService.ListActiveServices(c.Request.Context())
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "failed to list services")
		return
	}

	response := make([]serviceResponse, 0, len(services))
	for _, item := range services {
		response = append(response, serviceResponse{
			ID:              item.ID,
			Name:            item.Name,
			DurationMinutes: item.DurationMinutes,
			Price:           item.Price,
			IsActive:        item.IsActive,
			CreatedAt:       item.CreatedAt.Format(timeLayout),
			UpdatedAt:       item.UpdatedAt.Format(timeLayout),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"services": response,
	})
}

func (h *ServiceHandler) errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}
