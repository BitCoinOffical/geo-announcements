package api

import (
	"net/http"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	h         *handlers.IncidentHandler
	sh        *handlers.SystemHandler
	loc       *handlers.LocationHandler
	appApiKey string
	engine    *gin.Engine
}

func NewServer() *Server {
	engine := gin.New()
	return &Server{engine: engine}
}

func (s *Server) RegisterRoutes(h *handlers.IncidentHandler, sh *handlers.SystemHandler, loc *handlers.LocationHandler, appApiKey string) *Server {
	return &Server{h: h, sh: sh, loc: loc, appApiKey: appApiKey}
}

func (s *Server) Routes() {

	s.engine.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/api/v1")
	})

	rout := s.engine.Group("/api/v1")
	{
		rout.GET("/incidents/", middleware.CheckApiKey(s.appApiKey), s.h.GetIncidentsHandler)
		rout.GET("/incidents", middleware.CheckApiKey(s.appApiKey), s.h.GetIncidentByIDHandler)
		rout.POST("/incidents", middleware.CheckApiKey(s.appApiKey), s.h.CreateIncidentsHandler)
		rout.PUT("/incidents", middleware.CheckApiKey(s.appApiKey), s.h.UpdateIncidentsByIDHandler)
		rout.DELETE("/incidents", middleware.CheckApiKey(s.appApiKey), s.h.DeleteIncidentsByIDHandler)
		rout.GET("/incidents/stats", middleware.CheckApiKey(s.appApiKey), s.h.GetIncidentStatHandler)
		rout.GET("/system/health", middleware.CheckApiKey(s.appApiKey), s.sh.GetSystemHealth)
		rout.POST("/location/check", s.loc.CreateLocationHandler)
	}

}
func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}
