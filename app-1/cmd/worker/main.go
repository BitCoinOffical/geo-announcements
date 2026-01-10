package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	rdb "github.com/BitCoinOffical/geo-announcements/app-1/internal/adapters/secondary/redis"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/redis/go-redis/v9"
)

type WebhookWorker struct {
	rdb     *redis.Client
	Workers int
	Retry   int
	Key     string
}

func main() {
	cfg, err := config.NewLoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	rdb := rdb.NewWebhookRedis(&cfg.Redis)

	w := NewWebhookWorker(rdb, cfg.App.RETRY, cfg.App.WORKERS, cfg.App.QUEUE_KEY)
	w.WebhookWorker(context.Background())
}

func NewWebhookWorker(rdb *redis.Client, Retry, Workers int, Key string) *WebhookWorker {
	return &WebhookWorker{rdb: rdb, Retry: Retry, Workers: Workers, Key: Key}
}

func (w *WebhookWorker) WebhookWorker(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for range w.Workers {
		wg.Go(func() {
			for {
				result, err := w.rdb.BLPop(ctx, 0, w.Key).Result()
				if err != nil {
					log.Println("BLPop:", err)
					continue
				}

				playload := result[1]
				fmt.Println("playload:", playload)
				var webhook dto.WebHookDTO
				if err := json.Unmarshal([]byte(playload), &webhook); err != nil {
					log.Println("Unmarshal: ", err)
					continue
				}
				fmt.Println("webhook:", webhook)

				if err := sendWebhook(&webhook, playload); err == nil {
					log.Println("webhook sucess: ok")
					continue
				}
				webhook.RetryCount++
				if webhook.RetryCount > w.Retry {
					log.Println("send webhook failed")
					continue
				}

				time.AfterFunc(5*time.Second, func() {
					b, err := json.Marshal(webhook)
					if err != nil {
						log.Println(err)
					}
					w.rdb.RPush(ctx, "queue:webhooks", b)
				})

			}
		})
	}

	wg.Wait()
}

func sendWebhook(webhook *dto.WebHookDTO, playload string) error {
	_, err := http.Post(webhook.URL, "application/json", bytes.NewReader([]byte(playload)))
	return err
}
