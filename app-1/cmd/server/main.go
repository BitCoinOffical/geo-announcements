package main

import (
	"log"
	"net/http"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/adapters/secondary/migration"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/adapters/secondary/postgres"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/adapters/secondary/redis"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/domain/rules"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/handlers"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/queue"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.NewLoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var logger *zap.Logger
	switch cfg.App.DEBUG_LEVEL {
	case "Dev":
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatal(err)
		}
	case "Prod":
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
	migrationsDir := "./migrations"
	migration.RunMigrations(db, migrationsDir)

	rdb := redis.NewRedis(&cfg.Redis)
	qrd := redis.NewWebhookRedis(&cfg.Redis)

	h := handlers.NewIncidentHandler(services.NewIncidentService(repo.NewIncidentRepo(db), cache.NewIncidentCache(rdb), logger), logger, &cfg.App)
	sh := handlers.NewSystemHandler(db, rdb, logger)
	loc := handlers.NewLocationHandler(services.NewLocationService(repo.NewLocationRepo(db)), queue.NewWebHookQueue(qrd), logger, &cfg.App)

	r := gin.New()

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

	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/api/v1")
	})

	rout := r.Group("/api/v1")
	{
		rout.GET("/incidents/", middleware.CheckApiKey(&cfg.App), h.GetIncidentsHandler)
		rout.GET("/incidents", middleware.CheckApiKey(&cfg.App), h.GetIncidentByIDHandler)
		rout.POST("/incidents", middleware.CheckApiKey(&cfg.App), h.CreateIncidentsHandler)
		rout.PUT("/incidents", middleware.CheckApiKey(&cfg.App), h.UpdateIncidentsByIDHandler)
		rout.DELETE("/incidents", middleware.CheckApiKey(&cfg.App), h.DeleteIncidentsByIDHandler)
		rout.GET("/incidents/stats", middleware.CheckApiKey(&cfg.App), h.GetIncidentStatHandler)
		rout.GET("/system/health", middleware.CheckApiKey(&cfg.App), sh.GetSystemHealth)
		rout.POST("/location/check", loc.CreateLocationHandler)
	}

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
