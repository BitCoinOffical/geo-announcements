package handlers

import (
	"context"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
)

type incidentService interface {
	GetIncidents(ctx context.Context, page, limit int) ([]models.Incident, error)
	GetIncidentByID(ctx context.Context, id int) (*models.Incident, error)
	GetIncidentStat(ctx context.Context, cfg *config.AppConfig) (*models.UsersInDangerousZones, error)
	CreateIncidents(ctx context.Context, dto *dto.IncidentDTO) error
	UpdateIncidentsByID(ctx context.Context, dto *dto.IncidentDTO, id int) error
	DeleteIncidentsByID(ctx context.Context, id int) error
}

type locationService interface {
	CreateLocation(ctx context.Context, dto *dto.LocationDTO, userID string) ([]models.DangerousZones, error)
}

type Handlers struct {
	Incident incidentService
	Location locationService
}

func NewServices(Incident incidentService, Location locationService) *Handlers {
	return &Handlers{Incident: Incident, Location: Location}
}
