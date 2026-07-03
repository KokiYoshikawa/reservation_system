package handler

import (
	"net/http"

	"reservation-system/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": h.repo.GetRootMessage(),
	})
}

func (h *Handler) Health(c *gin.Context) {
	if err := h.repo.CheckHealth(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "database unavailable",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": h.repo.GetHealthStatus(),
	})
}
