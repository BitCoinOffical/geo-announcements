package services

import (
	"context"

	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/internal/interfaces/http/repo"
)

type LocationService struct {
	repo *repo.LocationRepo
}

func NewLocationService(repo *repo.LocationRepo) *LocationService {
	return &LocationService{repo: repo}
}

func (h *LocationService) CreateLocationService(ctx context.Context, dto *dto.LocationDTO, userID string) error {
	if err := h.repo.CreateLocationRepo(ctx, dto, userID); err != nil {
		return err
	}
	return nil
}
