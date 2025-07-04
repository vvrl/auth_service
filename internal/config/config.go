package config

import (
	"time"

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
		Driver      string        `yaml:"driver"`
		Host        string        `yaml:"host"`
		Port        int           `yaml:"port"`
		User        string        `yaml:"user"`
		Password    string        `yaml:"password"`
		Dbname      string        `yaml:"dbname"`
		MaxAttempts int           `yaml:"max_attempts"`
		Timeout     time.Duration `yaml:"timeout"`
	}
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("config read error: %v", err)
	}

	var cfg Config

	err := viper.Unmarshal(&cfg)
	if err != nil {
		logrus.Fatalf("parsing in struct error: %v", err)
	}

	return &cfg
}
