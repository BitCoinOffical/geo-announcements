package services

import (
	"context"
	"fmt"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/repo"
	"go.uber.org/zap"
)

const (
	topLimit       = 100
	maxPage        = 10
	defaultPageTTL = 30
	defaultByIdTTL = 5
)

type IncidentService struct {
	logger       *zap.Logger
	incidentRepo *repo.IncidentRepo
	cache        *cache.IncidentCache
}

func NewIncidentService(incidentRepo *repo.IncidentRepo, cache *cache.IncidentCache, logger *zap.Logger) *IncidentService {
	return &IncidentService{incidentRepo: incidentRepo, cache: cache, logger: logger}
}

func (h *IncidentService) GetIncidentsService(ctx context.Context, page, limit int) ([]models.Incident, error) {
	offset := (page - 1) * limit
	if page <= maxPage {
		all, err := h.cache.GetTop(ctx)
		if err != nil {
			return nil, fmt.Errorf("cache.GetTop: %w", err)
		}
		if len(all) > 0 {
			start := (page - 1) * limit
			if start >= len(all) {
				return nil, nil
			}

			end := start + limit
			if end > len(all) {
				end = len(all)
			}

			h.logger.Debug("received from top redis", zap.String("source", "redis"))
			return all[start:end], nil
		}

		rows, err := h.incidentRepo.GetTopRepo(ctx, topLimit)
		if err != nil {
			return nil, fmt.Errorf("incidentRepo.GetTopRepo: %w", err)
		}

		if err := h.cache.SetTop(ctx, rows, time.Minute); err != nil {
			return nil, fmt.Errorf("cache.SetTop: %w", err)
		}

		h.logger.Debug("received from top postgres and save top redis", zap.String("source", "postgres"))
		return rows, nil
	}

	rows, err := h.cache.GetAll(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("cache.GetAll: %w", err)
	}

	if rows != nil {
		h.logger.Debug("received from redis", zap.String("source", "redis"))
		return rows, nil
	}

	rows, err = h.incidentRepo.GetIncidentsRepo(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("incidentRepo.GetIncidentsRepo: %w", err)
	}

	if err := h.cache.Set(ctx, fmt.Sprintf("incident-page%d:-limit:%d", page, limit), rows, defaultPageTTL*time.Second); err != nil {
		return nil, fmt.Errorf("cache.Set: %w", err)
	}

	h.logger.Debug("received from postgres and saved in redis", zap.String("source", "postgres"))
	return rows, nil
}

func (h *IncidentService) GetIncidentByIDService(ctx context.Context, id int) (*models.Incident, error) {
	row, err := h.cache.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cache.Get %w", err)
	}
	if row != nil {
		h.logger.Debug("received from redis", zap.String("source", "redis"))
		return row, nil
	}

	row, err = h.incidentRepo.GetIncidentByIDRepo(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("incidentRepo.GetIncidentByIDRepo %w", err)
	}

	if err := h.cache.Set(ctx, fmt.Sprintf("incident:%d", id), row, defaultByIdTTL*time.Minute); err != nil {
		h.logger.Debug("saved in redis")
		return nil, fmt.Errorf("cache.Set %w", err)
	}
	h.logger.Debug("received from postgres", zap.String("source", "postgres"))
	return row, nil
}

func (h *IncidentService) GetIncidentStatService(ctx context.Context, cfg *config.AppConfig) (*models.UsersInDangerousZones, error) {
	fromTime := time.Now().Add(-time.Duration(cfg.STATS_TIME_WINDOW_MINUTES) * time.Minute)
	users, err := h.incidentRepo.GetIncidentStatRepo(ctx, &fromTime)
	if err != nil {
		return nil, fmt.Errorf("incidentRepo.GetIncidentStatRepo: %w", err)
	}
	return users, nil
}

func (h *IncidentService) CreateIncidentsService(ctx context.Context, dto *dto.IncidentDTO) error {
	if err := h.incidentRepo.CreateIncidentsRepo(ctx, dto); err != nil {
		return fmt.Errorf("incidentRepo.CreateIncidentsRepo: %w", err)
	}
	return nil
}

func (h *IncidentService) UpdateIncidentsByIDService(ctx context.Context, dto *dto.IncidentDTO, id int) error {

	if err := h.incidentRepo.UpdateIncidentsByIDRepo(ctx, dto, id); err != nil {
		return fmt.Errorf("incidentRepo.UpdateIncidentsByIDRepo: %w", err)
	}

	if err := h.cache.Del(ctx, id); err != nil {
		return fmt.Errorf("cache.Del: %w", err)
	}

	return nil
}

func (h *IncidentService) DeleteIncidentsByIDService(ctx context.Context, id int) error {
	if err := h.incidentRepo.DeleteIncidentsByIDRepo(ctx, id); err != nil {
		return fmt.Errorf("incidentRepo.DeleteIncidentsByIDRepo: %w", err)
	}

	if err := h.cache.Del(ctx, id); err != nil {
		return fmt.Errorf("cache.Del: %w", err)
	}

	return nil
}
