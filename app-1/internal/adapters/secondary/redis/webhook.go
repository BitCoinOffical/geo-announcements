package redis

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewWebhookRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("RDB_HOST"), os.Getenv("RDB_PORT")),
		Password: os.Getenv("RDB_PASSWORD"),
		DB:       1,
	})
	return rdb
}
