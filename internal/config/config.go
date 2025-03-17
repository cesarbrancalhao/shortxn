package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresURL     string  `envconfig:"POSTGRES_URL" default:"postgres://shortxn:shortxn123@localhost:5432/shortxn?sslmode=disable"`
	RedisAddr       string  `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	RabbitMQURL     string  `envconfig:"RABBITMQ_URL" default:"amqp://guest:guest@localhost:5672/"`
	ServerPort      string  `envconfig:"SERVER_PORT" default:":8080"`
	RateLimit       float64 `envconfig:"RATE_LIMIT" default:"100"`
	RateBurst       int     `envconfig:"RATE_BURST" default:"50"`
	MaxURLLength    int     `envconfig:"MAX_URL_LENGTH" default:"2048"`
	CacheExpiration int     `envconfig:"CACHE_EXPIRATION" default:"24"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
