package authservice

import (
	"auth_service/internal/api"
	"auth_service/internal/auth"
	"auth_service/internal/config"
	"auth_service/internal/storage"
	"context"

	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	config  *config.Config
	logger  *logrus.Logger
	router  *mux.Router
	store   *storage.Store
	service *auth.AuthService
	handler *api.Handler
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

	a.service = auth.NewService(a.store, a.logger)
	a.handler = api.NewHandler(a.service, a.store, a.logger)

	e := echo.New()
	api.CreateRouters(e, a.handler)

	e.Logger.Fatal(e.Start(":" + a.config.Server.Port))

	a.logger.Info("сервер запущен на порту: ", a.config.Server.Port)

	return nil
}
