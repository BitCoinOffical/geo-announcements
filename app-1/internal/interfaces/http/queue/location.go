package queue

import (
	"context"
	"encoding/json"

	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/dto"
	"github.com/redis/go-redis/v9"
)

type WebHookQueue struct {
	rdb *redis.Client
}

func NewWebHookQueue(rdb *redis.Client) *WebHookQueue {
	return &WebHookQueue{rdb: rdb}
}

func (w *WebHookQueue) EnqueueWebHook(ctx context.Context, webhook *dto.WebHookDTO) error {
	data, err := json.Marshal(webhook)
	if err != nil {
		return err
	}
	return w.rdb.RPush(ctx, "queue:webhooks", data).Err()
}
