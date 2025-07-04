package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	}

	Logger struct {
		Level    string `yaml:"level"`
		FileName string `yaml:"file_name"`
	}

	Database struct {
		Driver      string `yaml:"driver"`
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		Dbname      string `yaml:"dbname"`
		MaxAttempts int    `yaml:"maxAttempts"`
	}
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("ошибка чтения конфига: %v", err)
	}

	var cfg Config

	err := viper.Unmarshal(&cfg)
	if err != nil {
		logrus.Fatalf("неудачный парсинг конфига в структуру: %v", err)
	}

	return &cfg
}
