package services

import (
	"context"
	"fmt"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/repo"
)

type LocationService struct {
	repo *repo.LocationRepo
}

func NewLocationService(repo *repo.LocationRepo) *LocationService {
	return &LocationService{repo: repo}
}

func (h *LocationService) CreateLocation(ctx context.Context, dto *dto.LocationDTO, userID string) ([]models.DangerousZones, error) {
	zones, err := h.repo.CreateLocation(ctx, dto, userID)
	if err != nil {
		return nil, fmt.Errorf("repo.CreateLocationRepo: %w", err)
	}
	return zones, nil
}
