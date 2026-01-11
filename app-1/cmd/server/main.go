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
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/queue"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	Dev           = "Dev"
	Prod          = "Prod"
	migrationsDir = "./migrations"
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

	h := handlers.NewIncidentHandler(services.NewIncidentService(repo.NewIncidentRepo(db), cache.NewIncidentCache(cacheRdb), logger), logger, &cfg.App)
	sh := handlers.NewSystemHandler(db, cacheRdb, logger)
	loc := handlers.NewLocationHandler(services.NewLocationService(repo.NewLocationRepo(db)), queue.NewWebHookQueue(queueRdb), logger, &cfg.App)

	api.NewRouters(h, sh, loc, cfg.App.ApiKey).Routers()

}
