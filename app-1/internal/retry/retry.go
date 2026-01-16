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

type Retry struct {
	logger *zap.Logger
	queue  *queue.WebHookQueue
	cfg    *config.AppConfig
}

func NewRetry(logger *zap.Logger, queue *queue.WebHookQueue, cfg *config.AppConfig) *Retry {
	return &Retry{logger: logger, queue: queue, cfg: cfg}

}

func (r *Retry) Retry(ctx context.Context, webhook *dto.WebHookDTO) {
	for i := range maxAttempts {
		if i == maxAttempts {
			r.logger.Error("exceeded the number of attempts")
			break
		}

		err := r.queue.EnqueueWebHook(ctx, webhook, r.cfg.QueueKey)
		if err == nil {
			return
		}
		r.logger.Error("enqueue webhook error", zap.Error(err))
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			continue
		}
	}
}
