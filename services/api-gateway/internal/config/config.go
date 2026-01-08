package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer    `yaml:"http"`
	Services      `yaml:"services"`
	Auth          `yaml:"auth"`
	LoggingConfig `yaml:"logging"`
}

type HTTPServer struct {
	Port        int           `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	ReadTimeout time.Duration `yaml:"read_timeout" env-default:"60s"`
}

type Services struct {
	AuthServiceAddr  string `yaml:"auth_service" env:"AUTH_SERVICE_ADDR" env-required:"true"`
	BoardServiceAddr string `yaml:"board_service" env:"BOARD_SERVICE_ADDR" env-required:"true"`
}

type Auth struct {
	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
}

type LoggingConfig struct {
	Level string `yaml:"Level"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	env := "dev"
	if e := os.Getenv("APP_ENV"); e != "" {
		env = e
	}

	configPath := fmt.Sprintf("config/config.%s.yaml", env)

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, errors.New("config file not found")
	}
	return &cfg, nil
}
