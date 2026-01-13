package services

import (
	"context"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
)

type incidentRepository interface {
	GetTop(ctx context.Context, limit int) ([]models.Incident, error)
	GetIncidents(ctx context.Context, page, limit int) ([]models.Incident, error)
	GetIncidentByID(ctx context.Context, id int) (*models.Incident, error)
	GetIncidentStat(ctx context.Context, fromTime *time.Time) (*models.UsersInDangerousZones, error)
	CreateIncidents(ctx context.Context, dto *dto.IncidentDTO) error
	UpdateIncidentsByID(ctx context.Context, dto *dto.IncidentDTO, id int) error
	DeleteIncidentsByID(ctx context.Context, id int) error
}

type incidentCache interface {
	GetTop(ctx context.Context) ([]models.Incident, error)
	SetTop(ctx context.Context, value any, ttl time.Duration) error
	GetAll(ctx context.Context, page int, limit int) ([]models.Incident, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, id int) (*models.Incident, error)
	Del(ctx context.Context, id int) error
}

type locationRepository interface {
	CreateLocation(ctx context.Context, dto *dto.LocationDTO, userID string) ([]models.DangerousZones, error)
}

type Service struct {
	IncidentRepo  incidentRepository
	LocationRepo  locationRepository
	IncidentCache incidentCache
}

func NewRepos(incidentRepo incidentRepository, locationRepo locationRepository, incidentCache incidentCache) *Service {
	return &Service{
		IncidentRepo:  incidentRepo,
		LocationRepo:  locationRepo,
		IncidentCache: incidentCache,
	}
}
