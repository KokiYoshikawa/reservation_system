package router

import (
	"reservation-system/backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func New(h *handler.Handler) (*gin.Engine, error) {
	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	router.GET("/", h.Root)

	api := router.Group("/api/v1")
	api.GET("/health", h.Health)

	return router, nil
}
