package authservice

import (
	"auth_service/internal/config"
	"auth_service/internal/handlers"
	"auth_service/internal/storage"
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	config *config.Config
	logger *logrus.Logger
	router *mux.Router
	store  *storage.Store
}

func NewApp(config *config.Config, logger *logrus.Logger) *AuthService {
	return &AuthService{
		config: config,
		logger: logger,
		router: mux.NewRouter(),
		store:  storage.NewStore(),
	}
}

func (a *AuthService) Run() error {

	if err := a.store.Open(context.Background(), a.config, a.logger); err != nil {
		a.logger.Fatal("ошибка при подключении к базе данных: ", err)
	}
	a.logger.Info("успешное подключение к базе данных")
	defer a.store.Close()

	handlers.CreateRouters(a.router)

	a.logger.Info("сервер запущен на порту: ", a.config.Server.Port)

	return http.ListenAndServe(":"+a.config.Server.Port, a.router)
}
