package services

import (
	"context"
	"fmt"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
)

type LocationService struct {
	repo locationRepository
}

func NewLocationService(repo locationRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (h *LocationService) CreateLocation(ctx context.Context, dto *dto.LocationDTO, userID string) error {
	err := h.repo.CreateLocation(ctx, dto, userID)
	if err != nil {
		return fmt.Errorf("repo.CreateLocationRepo: %w", err)
	}
	return nil
}

func (h *LocationService) GetDangerZones(ctx context.Context, dto *dto.LocationDTO, userID string) ([]models.DangerousZones, error) {
	zones, err := h.repo.GetDangerZones(ctx, dto, userID)
	if err != nil {
		return nil, err
	}
	return zones, nil
}
