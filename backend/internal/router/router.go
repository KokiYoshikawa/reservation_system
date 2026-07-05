package router

import (
	"reservation-system/backend/internal/handler"
	"reservation-system/backend/internal/middleware"
	"reservation-system/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func New(h *handler.Handler, authHandler *handler.AuthHandler, authService service.AuthService, userHandler *handler.UserHandler) (*gin.Engine, error) {
	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	router.GET("/", h.Root)

	api := router.Group("/api/v1")
	api.GET("/health", h.Health)
	auth := api.Group("/auth")
	authProtected := auth.Group("")
	authProtected.Use(middleware.NewAuthMiddleware(authService))
	auth.POST("/register", userHandler.CreateUser)
	auth.POST("/login", authHandler.Login)
	authProtected.POST("/logout", authHandler.Logout)
	authProtected.GET("/me", authHandler.Me)

	return router, nil
}
