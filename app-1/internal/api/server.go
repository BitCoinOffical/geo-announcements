package api

import (
	"log"
	"net/http"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/handlers"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

type Routers struct {
	h         *handlers.IncidentHandler
	sh        *handlers.SystemHandler
	loc       *handlers.LocationHandler
	appApiKey string
}

func NewRouters(h *handlers.IncidentHandler, sh *handlers.SystemHandler, loc *handlers.LocationHandler, appApiKey string) *Routers {
	return &Routers{h: h, sh: sh, loc: loc, appApiKey: appApiKey}
}

func (r *Routers) Routers() {

	g := gin.New()

	g.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/api/v1")
	})

	rout := g.Group("/api/v1")
	{
		rout.GET("/incidents/", middleware.CheckApiKey(r.appApiKey), r.h.GetIncidentsHandler)
		rout.GET("/incidents", middleware.CheckApiKey(r.appApiKey), r.h.GetIncidentByIDHandler)
		rout.POST("/incidents", middleware.CheckApiKey(r.appApiKey), r.h.CreateIncidentsHandler)
		rout.PUT("/incidents", middleware.CheckApiKey(r.appApiKey), r.h.UpdateIncidentsByIDHandler)
		rout.DELETE("/incidents", middleware.CheckApiKey(r.appApiKey), r.h.DeleteIncidentsByIDHandler)
		rout.GET("/incidents/stats", middleware.CheckApiKey(r.appApiKey), r.h.GetIncidentStatHandler)
		rout.GET("/system/health", middleware.CheckApiKey(r.appApiKey), r.sh.GetSystemHealth)
		rout.POST("/location/check", r.loc.CreateLocationHandler)
	}

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
