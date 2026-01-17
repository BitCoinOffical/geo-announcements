package handlers

import (
	"context"
	"database/sql"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/queue"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/retry"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type incidentService interface {
	GetIncidents(ctx context.Context, page, limit int) ([]models.Incident, error)
	GetIncidentByID(ctx context.Context, id int) (*models.Incident, error)
	GetIncidentStat(ctx context.Context, StatsTimeWindowMinutes int) (*models.UsersInDangerousZones, error)
	CreateIncidents(ctx context.Context, dto *dto.IncidentDTO) error
	UpdateZones(ctx context.Context, dto *dto.IncidentDTO) error
	UpdateIncidentsByID(ctx context.Context, dto *dto.IncidentDTO, id int) error
	DeleteIncidentsByID(ctx context.Context, id int) error
}

type locationService interface {
	CreateLocation(ctx context.Context, dto *dto.LocationDTO, userID string) error
	GetDangerZones(ctx context.Context, dto *dto.LocationDTO) ([]models.DangerousZones, error)
}

type webhookRetry interface {
	Retry(ctx context.Context, webhook *dto.WebHookDTO)
}

type DBPinger interface {
	Ping() error
}

type RedisPinger interface {
	Ping(ctx context.Context) *redis.StatusCmd
}

type Services struct {
	Incident incidentService
	Location locationService
	Retry    webhookRetry
}

func NewServices(db *sql.DB, queueRdb, cacheRdb *redis.Client, logger *zap.Logger, cfg *config.Config) *Services {
	incidentRepo := repo.NewIncidentRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	incidentCache := cache.NewIncidentCache(cacheRdb)
	queue := queue.NewWebHookQueue(queueRdb)

	retry := retry.NewRetry(logger, queue, &cfg.App)
	incSrv := services.NewIncidentService(incidentRepo, incidentCache, logger)
	locSrv := services.NewLocationService(locationRepo)

	return &Services{Incident: incSrv, Location: locSrv, Retry: retry}
}

type Handlers struct {
	Incident *IncidentHandler
	System   *SystemHandler
	Location *LocationHandler
}

func NewHandlers(svcs *Services, db *sql.DB, cacheRdb *redis.Client, logger *zap.Logger, cfg *config.Config) *Handlers {
	h := NewIncidentHandler(svcs.Incident, logger, &cfg.App)
	sh := NewSystemHandler(db, cacheRdb, logger)
	loc := NewLocationHandler(svcs.Location, logger, &cfg.App, svcs.Retry)
	return &Handlers{Incident: h, System: sh, Location: loc}
}
