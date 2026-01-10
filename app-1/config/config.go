package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	App      AppConfig
}

type PostgresConfig struct {
	DB_HOST     string `env:"DB_HOST"`
	DB_PORT     string `env:"DB_PORT"`
	DB_USER     string `env:"DB_USER"`
	DB_PASSWORD string `env:"DB_PASSWORD"`
	DB_NAME     string `env:"DB_NAME"`
}

type RedisConfig struct {
	RDB_HOST     string `env:"RDB_HOST"`
	RDB_PORT     string `env:"RDB_PORT"`
	RDB_PASSWORD string `env:"RDB_PASSWORD"`
}

type AppConfig struct {
	API_KEY                   string `env:"API_KEY"`
	QUEUE_KEY                 string `env:"QUEUE_KEY"`
	STATS_TIME_WINDOW_MINUTES int    `env:"STATS_TIME_WINDOW_MINUTES"`
	WEBHOOK_URL               string `env:"WEBHOOK_URL"`
	WORKERS                   int    `env:"WORKERS"`
	RETRY                     int    `env:"RETRY"`
	DEBUG_LEVEL               string `env:"DEBUG_LEVEL"`
}

func NewLoadConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg.Postgres); err != nil {
		return nil, err
	}
	if err := env.Parse(&cfg.Redis); err != nil {
		return nil, err
	}
	if err := env.Parse(&cfg.App); err != nil {
		return nil, err
	}
	return &cfg, nil
}
