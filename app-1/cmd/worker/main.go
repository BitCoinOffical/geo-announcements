package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	redi "github.com/BitCoinOffical/geo-announcements/internal/adapters/secondary/redis"
	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/dto"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("no .env file found")
	}
	workers, err := strconv.Atoi(os.Getenv("WORKERS"))
	if err != nil {
		workers = 1
		log.Println(err)
	}
	retry, err := strconv.Atoi(os.Getenv("RETRY"))
	if err != nil {
		retry = 5
		log.Println("error:", err)
	}
	rdb := redi.NewWebhookRedis()
	WebhookWorker(context.Background(), rdb, workers, retry)
}
func WebhookWorker(ctx context.Context, rdb *redis.Client, workers int, retry int) {
	wg := &sync.WaitGroup{}
	for range workers {
		wg.Go(func() {
			for {
				result, err := rdb.BLPop(ctx, 0, "queue:webhooks").Result()
				if err != nil {
					log.Println("BLPop>", err)
					continue
				}

				playload := result[1]
				fmt.Println("playload>>> ", playload)
				var webhook dto.WebHookDTO
				if err := json.Unmarshal([]byte(playload), &webhook); err != nil {
					log.Println("Unmarshal>", err)
					continue
				}
				fmt.Println("webhook>>> ", webhook)

				if err := sendWebhook(&webhook, playload); err == nil {
					log.Println("webhook sucess: ok")
					continue
				}
				webhook.RetryCount++
				if webhook.RetryCount > retry {
					log.Println("send webhook failed")
					continue
				}

				time.AfterFunc(5*time.Second, func() {
					b, err := json.Marshal(webhook)
					if err != nil {
						log.Println(err)
					}
					rdb.RPush(ctx, "queue:webhooks", b)
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
