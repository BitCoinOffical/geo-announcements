package main

import (
	"log"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/adapters/secondary/migration"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/adapters/secondary/postgres"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/adapters/secondary/redis"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/domain/rules"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	Dev           = "Dev"
	Prod          = "Prod"
	migrationsDir = "./migrations"
)

// основной сервис
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

	db, err := postgres.NewPostgres(&cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	migration.RunMigrations(db, migrationsDir)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("lat", rules.ValidateLat); err != nil {
			log.Printf("error lat validate: %v", err)
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("lon", rules.ValidateLon); err != nil {
			log.Printf("error lat validate: %v", err)
		}
	}

	cacheRdb := redis.NewRedis(&cfg.Redis)
	queueRdb := redis.NewWebhookRedis(&cfg.Redis)

	services := handlers.NewServices(db, queueRdb, cacheRdb, logger, cfg)
	hanlrs := handlers.NewHandlers(services, db, cacheRdb, logger, cfg)
	serv := api.NewServer(hanlrs)
	err = serv.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
