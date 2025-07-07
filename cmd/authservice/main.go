package main

import (
	"auth_service/internal/authservice"
	"auth_service/internal/config"
	"auth_service/internal/logger"
	"log"
)

func main() {
	cfg := config.NewConfig()
	logger, err := logger.NewLogger(cfg)

	if err != nil {
		log.Fatal("ошибка при создании логера:", err)
	}

	authApp := authservice.NewApp(cfg, logger)

	if err := authApp.Run(); err != nil {
		logger.Fatal("ошибка при запуске сервера: ", err)
	}
}
