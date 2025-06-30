package main

import (
	"auth_service/internal/authservice"
	"log"
)

func main() {
	config := authservice.NewConfig()
	logger, err := authservice.NewLogger(config)

	if err != nil {
		log.Fatal("Ошибка при создании логера:", err)
	}

	authService := authservice.NewApp(config, logger)

	if err := authService.Run(); err != nil {
		logger.Fatal("Ошибка при запуске сервера: ", err)
	}
}
