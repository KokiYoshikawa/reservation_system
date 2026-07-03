package main

import (
	"log"
	"os"

	"reservation-system/backend/internal/handler"
	"reservation-system/backend/internal/repository"

	internalrouter "reservation-system/backend/internal/router"
)

func main() {
	repo := repository.New()
	h := handler.New(repo)
	router, err := internalrouter.New(h)
	if err != nil {
		log.Fatalf("failed to configure trusted proxies: %v", err)
	}

	addr := ":" + getEnv("PORT", "8080")
	log.Printf("starting gin server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start gin server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
