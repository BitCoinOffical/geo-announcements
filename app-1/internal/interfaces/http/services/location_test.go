package services_test

import (
	"errors"
	"testing"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services"
	mock_services "github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// набор тестов для CreateLocation
func TestCreateLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMocklocationRepository(ctrl)
	services := services.NewLocationService(repo)
	dto := &dto.LocationDTO{}
	repo.EXPECT().CreateLocation(gomock.Any(), dto, "id").Return(nil)
	err := services.CreateLocation(t.Context(), dto, "id")

	assert.NoError(t, err)
}

func TestCreateLocation_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMocklocationRepository(ctrl)
	services := services.NewLocationService(repo)
	dto := &dto.LocationDTO{}
	repo.EXPECT().CreateLocation(gomock.Any(), dto, "id").Return(errors.New("db is down")) //имитируем ошибку базы данных
	err := services.CreateLocation(t.Context(), dto, "id")

	assert.Error(t, err)
}

func TestGetDangerZones(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMocklocationRepository(ctrl)
	services := services.NewLocationService(repo)
	dto := &dto.LocationDTO{}
	repo.EXPECT().GetDangerZones(gomock.Any(), dto, "id").Return([]models.DangerousZones{}, nil) //имитируем ошибку базы данных
	zones, err := services.GetDangerZones(t.Context(), dto, "id")

	assert.Len(t, zones, 0)
	assert.NoError(t, err)
}
func TestGetDangerZones_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMocklocationRepository(ctrl)
	services := services.NewLocationService(repo)
	dto := &dto.LocationDTO{}
	repo.EXPECT().GetDangerZones(gomock.Any(), dto, "id").Return(nil, errors.New("db is down")) //имитируем ошибку базы данных
	_, err := services.GetDangerZones(t.Context(), dto, "id")

	assert.Error(t, err)

}
