package api

import (
	"net/http"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	h      *handlers.Handlers
}

func NewServer(h *handlers.Handlers) *Server {
	return &Server{h: h, engine: gin.New()}
}

func (s *Server) Run(cfg *config.Config) error {

	s.engine.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/api/v1")
	})

	appApiKey := cfg.App.ApiKey
	rout := s.engine.Group("/api/v1")
	{
		rout.GET("/incidents/", middleware.CheckApiKey(appApiKey), s.h.Incident.GetIncidentsHandler)
		rout.GET("/incidents", middleware.CheckApiKey(appApiKey), s.h.Incident.GetIncidentByIDHandler)
		rout.POST("/incidents", middleware.CheckApiKey(appApiKey), s.h.Incident.CreateIncidentsHandler)
		rout.PUT("/incidents", middleware.CheckApiKey(appApiKey), s.h.Incident.UpdateIncidentsByIDHandler)
		rout.DELETE("/incidents", middleware.CheckApiKey(appApiKey), s.h.Incident.DeleteIncidentsByIDHandler)
		rout.GET("/incidents/stats", middleware.CheckApiKey(appApiKey), s.h.Incident.GetIncidentStatHandler)
		rout.GET("/system/health", middleware.CheckApiKey(appApiKey), s.h.System.GetSystemHealth)
		rout.POST("/location/check", s.h.Location.CreateLocationHandler)
	}
	return s.engine.Run()
}
