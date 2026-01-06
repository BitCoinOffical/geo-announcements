package main

import (
	"log"
	"net/http"

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
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	db, err := postgres.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewRedis()
	qrd := redis.NewWebhookRedis()

	h := handlers.NewIncidentHandler(services.NewIncidentService(repo.NewIncidentRepo(db), cache.NewIncidentCache(rdb)))
	sh := handlers.NewSystemHandler(db, rdb)
	loc := handlers.NewLocationHandler(services.NewLocationService(repo.NewLocationRepo(db)), queue.NewWebHookQueue(qrd))

	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("lat", rules.ValidateLat)
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("lon", rules.ValidateLon)
	}

	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/api/v1")
	})

	rout := r.Group("/api/v1")
	{
		rout.GET("/incidents/", middleware.CheckApiKey(), h.GetIncidentsHandler)
		rout.GET("/incidents", middleware.CheckApiKey(), h.GetIncidentByIDHandler)
		rout.POST("/incidents", middleware.CheckApiKey(), h.CreateIncidentsHandler)
		rout.PUT("/incidents", middleware.CheckApiKey(), h.UpdateIncidentsByIDHandler)
		rout.DELETE("/incidents", middleware.CheckApiKey(), h.DeleteIncidentsByIDHandler)
		rout.GET("/incidents/stats", middleware.CheckApiKey(), h.GetIncidentStatHandler)
		rout.GET("/system/health", middleware.CheckApiKey(), sh.GetSystemHealth)
		rout.POST("/location/check", loc.CreateLocationHandler)
	}

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
