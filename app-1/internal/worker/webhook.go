package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/queue"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/retry"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type WebhookWorker struct {
	cfg    *config.AppConfig
	queue  *queue.WebHookQueue
	logger *zap.Logger
	rdb    *redis.Client
}

func NewWebhookWorker(rdb *redis.Client, logger *zap.Logger, cfg *config.AppConfig, queue *queue.WebHookQueue) *WebhookWorker {
	return &WebhookWorker{rdb: rdb, logger: logger, cfg: cfg, queue: queue}
}

func (w *WebhookWorker) WebhookWorker(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for range w.cfg.WorkersCount {
		wg.Go(func() {
			for {

				result, err := w.rdb.BLPop(ctx, 0, w.cfg.QueueKey).Result()
				if err != nil {
					if errors.Is(err, ctx.Err()) {
						return
					}
					w.logger.Error("BLPop error", zap.Error(err))
					continue
				}

				playload := result[1]
				w.logger.Debug("playload", zap.String("values", playload))
				var webhook dto.WebHookDTO
				if err := json.Unmarshal([]byte(playload), &webhook); err != nil {
					w.logger.Error("Unmarshal error", zap.Error(err))
					continue
				}

				w.logger.Debug("webhook", zap.Any("webhook", webhook))

				if err := sendWebhook(&webhook, playload); err == nil {
					w.logger.Info("webhook sucess: ok")
					continue
				}
				webhook.RetryCount++
				if webhook.RetryCount > w.cfg.SendWebhookRetry {
					w.logger.Warn("send webhook failed", zap.Int("exceeded the number of attempts", webhook.RetryCount))
					continue
				}

				retry.Retry(ctx, &webhook, w.logger, w.queue, w.cfg)

			}
		})
	}

	wg.Wait()
}

func sendWebhook(webhook *dto.WebHookDTO, playload string) error {
	_, err := http.Post(webhook.URL, "application/json", bytes.NewReader([]byte(playload)))
	return err
}
