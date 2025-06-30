package authservice

import (
	"auth_service/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

func NewApp(config *Config, logger *logrus.Logger) *AuthService {
	return &AuthService{
		config: config,
		logger: logger,
		router: mux.NewRouter(),
	}
}

func (a *AuthService) Run() error {

	// ПОДНЯТИЕ БД

	handlers.CreateRouters(a.router)

	a.logger.Info("starting server on port: ", a.config.Server.Port)

	return http.ListenAndServe(":"+a.config.Server.Port, a.router)
}
