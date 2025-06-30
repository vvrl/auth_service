package authservice

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(config *Config) (*logrus.Logger, error) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.Logger.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	// Set output file if specified
	// if config.Logger.FileName != "" {
	// 	file, err := os.OpenFile(config.Logger.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	logger.SetOutput(file)
	// } else {
	// 	logger.SetOutput(os.Stdout)
	// }

	return logger, nil
}
