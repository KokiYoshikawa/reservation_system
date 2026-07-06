package main

import (
	"log"

	"reservation-system/backend/internal/config"
	"reservation-system/backend/internal/database"
	"reservation-system/backend/internal/handler"
	"reservation-system/backend/internal/repository"
	"reservation-system/backend/internal/service"

	internalrouter "reservation-system/backend/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	redisClient, err := database.NewRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	repo := repository.New(db, redisClient, cfg.SessionTTL)
	h := handler.New(repo)
	authService := service.NewAuthService(repo, cfg)
	authHandler := handler.NewAuthHandler(authService)
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)
	serviceService := service.NewServiceService(repo)
	serviceHandler := handler.NewServiceHandler(serviceService)
	reservationSlotService := service.NewReservationSlotService(repo)
	reservationSlotHandler := handler.NewReservationSlotHandler(reservationSlotService)
	router, err := internalrouter.New(h, authHandler, authService, userHandler, serviceHandler, reservationSlotHandler)
	if err != nil {
		log.Fatalf("failed to configure trusted proxies: %v", err)
	}

	addr := ":" + cfg.Port
	log.Printf("starting gin server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start gin server: %v", err)
	}
}
