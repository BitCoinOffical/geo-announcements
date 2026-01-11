package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	rdb "github.com/BitCoinOffical/geo-announcements/app-1/internal/adapters/secondary/redis"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/queue"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/worker"
	"go.uber.org/zap"
)

const (
	Dev  = "Dev"
	Prod = "Prod"
)

func main() {
	cfg, err := config.NewLoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	var logger *zap.Logger
	switch cfg.App.DebugLevel {
	case Dev:
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatal(err)
		}
	case Prod:
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("incorrect debug value")
	}

	rdb := rdb.NewWebhookRedis(&cfg.Redis)
	queue := queue.NewWebHookQueue(rdb)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	w := worker.NewWebhookWorker(rdb, logger, &cfg.App, queue)
	w.WebhookWorker(ctx)
}
