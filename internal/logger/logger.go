package logger

import (
	"auth_service/internal/config"

	"github.com/sirupsen/logrus"
)

func NewLogger(config *config.Config) (*logrus.Logger, error) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.Logger.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	return logger, nil
}
