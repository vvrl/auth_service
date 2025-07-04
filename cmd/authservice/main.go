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
		log.Fatal("Ошибка при создании логера:", err)
	}

	authService := authservice.NewApp(cfg, logger)

	logger.Info("Сервис создан")

	if err := authService.Run(); err != nil {
		logger.Fatal("Ошибка при запуске сервера: ", err)
	}
}
