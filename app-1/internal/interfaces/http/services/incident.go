package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/models"
	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/repo"
)

type IncidentService struct {
	incidentRepo *repo.IncidentRepo
	cache        *cache.IncidentCache
}

func NewIncidentService(incidentRepo *repo.IncidentRepo, cache *cache.IncidentCache) *IncidentService {
	return &IncidentService{incidentRepo: incidentRepo, cache: cache}
}

func (h *IncidentService) GetIncidentsService(ctx context.Context, page, limit int) ([]models.Incident, error) {
	offset := (page - 1) * limit

	key := fmt.Sprintf("incident-page%d:-limit:%d-offset:%d", page, limit, offset)
	rows, err := h.cache.GetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	if rows != nil {
		log.Println("получено с редис")
		return rows, nil
	}

	rows, err = h.incidentRepo.GetIncidentsRepo(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	if err := h.cache.Set(ctx, key, rows, 30*time.Second); err != nil {
		return nil, err
	}
	log.Println("save с редис")
	return rows, err
}

func (h *IncidentService) GetIncidentByIDService(ctx context.Context, id int) (*models.Incident, error) {
	key := fmt.Sprintf("incident:%d", id)
	row, err := h.cache.Get(ctx, key)
	if err != nil {

		return nil, err
	}
	if row != nil {
		log.Println("получено с редис", row)
		return row, nil
	}

	row, err = h.incidentRepo.GetIncidentByIDRepo(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := h.cache.Set(ctx, key, row, 5*time.Minute); err != nil {
		return nil, err
	}
	log.Println("save с редис", row)

	return row, nil
}

func (h *IncidentService) GetIncidentStatService(ctx context.Context) (*int, error) {
	windowMinutes, err := strconv.Atoi(os.Getenv("STATS_TIME_WINDOW_MINUTES"))
	if err != nil {
		return nil, err
	}
	fromTime := time.Now().Add(-time.Duration(windowMinutes) * time.Minute)
	count, err := h.incidentRepo.GetIncidentStatRepo(ctx, &fromTime)
	if err != nil {
		return nil, err
	}
	return count, nil
}

func (h *IncidentService) CreateIncidentsService(ctx context.Context, dto *dto.IncidentDTO) error {
	if err := h.incidentRepo.CreateIncidentsRepo(ctx, dto); err != nil {
		return err
	}
	return nil
}

func (h *IncidentService) UpdateIncidentsByIDService(ctx context.Context, dto *dto.IncidentDTO, id int) error {

	if err := h.incidentRepo.UpdateIncidentsByIDRepo(ctx, dto, id); err != nil {
		return err
	}
	key := fmt.Sprintf("incident:%d", id)
	if err := h.cache.Del(ctx, key); err != nil {
		return err
	}
	return nil
}

func (h *IncidentService) DeleteIncidentsByIDService(ctx context.Context, id int) error {
	if err := h.incidentRepo.DeleteIncidentsByIDRepo(ctx, id); err != nil {
		return err
	}
	key := fmt.Sprintf("incident:%d", id)
	if err := h.cache.Del(ctx, key); err != nil {
		return err
	}
	return nil
}
