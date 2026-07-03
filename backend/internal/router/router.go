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
	auth := api.Group("/auth")
	auth.POST("/login", h.Login)
	auth.GET("/me", h.Me)
	auth.DELETE("/logout", h.Logout)

	return router, nil
}
