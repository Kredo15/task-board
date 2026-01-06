package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
	Auth     AuthConfig     `yaml:"auth"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type AppConfig struct {
	Name    string `yaml:"name" env-default:"auth-service"`
	Version string `yaml:"version"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"50051"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type PostgresConfig struct {
	URL         string `yaml:"url" env:"DB_URL" env-required:"true"`
	MaxPoolSize int    `yaml:"max_pool_size" env-default:"10"`
}

type RedisConfig struct {
	Host string `yaml:"host" env:"REDIS_HOST" env-default:"localhost:6379"`
	DB   int    `yaml:"db" env-default:"0"`
}

type AuthConfig struct {
	JWTSecret  string        `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
	AccessTTL  time.Duration `yaml:"access_ttl" env-default:"15m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-default:"720h"`
	BcryptCost int           `yaml:"bcrypt_cost" env-default:"12"`
}

type LoggingConfig struct {
	Level string `yaml:"Level"`
}

// MustLoad загружает конфиг или завершает программу при ошибке
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
