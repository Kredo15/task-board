package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPC     GRPCConfig     `yaml:"grpc"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
	Logging  LoggingConfig  `yaml:"logging"`
	Auth     AuthConfig     `yaml:"auth"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env:"GRPC_PORT" env-default:"50051"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

func (g GRPCConfig) Address() string {
	return fmt.Sprintf(":%d", g.Port)
}

type PostgresConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	DBname   string `yaml:"DBname"`

	MaxConns          int           `yaml:"DB_MAX_CONNS"`
	MinConns          int           `yaml:"DB_MIN_CONNS"`
	MaxConnLifetime   time.Duration `yaml:"DB_MAX_CONN_LIFETIME"`
	MaxConnIdleTime   time.Duration `yaml:"DB_MAX_CONN_IDLE_TIME"`
	HealthCheckPeriod time.Duration `yaml:"DB_HEALTH_CHECK_PERIOD"`
}

type AuthConfig struct {
	JWTSecret  string        `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
	AccessTTL  time.Duration `yaml:"access_ttl" env-default:"15m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-default:"720h"`
	BcryptCost int           `yaml:"bcrypt_cost" env-default:"12"`
}

func (p *PostgresConfig) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DBname, "disable",
	)
}

type RedisConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
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
