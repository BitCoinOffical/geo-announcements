package redis

import (
	"fmt"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/redis/go-redis/v9"
)

const (
	numWebHookDB = 1
)

func NewWebhookRedis(cfg *config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RDBHost, cfg.RDBPort),
		Password: cfg.RDBPassword,
		DB:       numWebHookDB,
	})
	return rdb
}
