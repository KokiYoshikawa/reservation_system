package router

import (
	"reservation-system/backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func New(h *handler.Handler, authHandler *handler.AuthHandler, userHandler *handler.UserHandler) (*gin.Engine, error) {
	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	router.GET("/", h.Root)

	api := router.Group("/api/v1")
	api.GET("/health", h.Health)
	auth := api.Group("/auth")
	auth.POST("/register", userHandler.CreateUser)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/me", authHandler.Me)
	auth.DELETE("/logout", authHandler.Logout)

	return router, nil
}
