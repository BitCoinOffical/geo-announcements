package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Env string

const (
	EnvProd = "prod"
	EnvDev  = "dev"
)

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	App      AppConfig
}

type PostgresConfig struct {
	DBHost     string `env:"DB_HOST"`
	DBport     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
}

type RedisConfig struct {
	RDBHost     string `env:"RDB_HOST"`
	RDBPort     string `env:"RDB_PORT"`
	RDBPassword string `env:"RDB_PASSWORD"`
}

type AppConfig struct {
	ApiKey                 string `env:"API_KEY"`
	QueueKey               string `env:"QUEUE_KEY"`
	StatsTimeWindowMinutes int    `env:"STATS_TIME_WINDOW_MINUTES"`
	WebhookUrl             string `env:"WEBHOOK_URL"`
	WorkersCount           int    `env:"WORKERS_COUNT"`
	SendWebhookRetry       int    `env:"SEND_WEBHOOK_RETRY"`
	DebugLevel             string `env:"DEBUG_LEVEL"`
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
	var env Env = Env(cfg.App.DebugLevel)
	if env != EnvProd && env != EnvDev {
		return nil, fmt.Errorf("incorrect debug level: %s", env)
	}
	return &cfg, nil
}
