package retry

import (
	"context"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/queue"
	"go.uber.org/zap"
)

const (
	maxAttempts = 5
)

func Retry(ctx context.Context, webhook *dto.WebHookDTO, logger *zap.Logger, queue *queue.WebHookQueue, cfg *config.AppConfig) {
	for i := range maxAttempts {
		if i == maxAttempts {
			logger.Error("exceeded the number of attempts")
			break
		}

		err := queue.EnqueueWebHook(ctx, webhook, cfg.QueueKey)
		if err == nil {
			return
		}
		logger.Error("enqueue webhook error", zap.Error(err))
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			continue
		}
	}
}
