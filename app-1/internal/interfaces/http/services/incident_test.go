package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services"
	mock_services "github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/services/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	id                     = 0
	page                   = 1
	limit                  = 10
	deafaultTTL            = 5
	topLimit               = 100
	key                    = "key"
	StatsTimeWindowMinutes = 5
)

// набор тестов для Top
func TestTop_FromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	service := services.NewIncidentService(repo, cache, logger)

	cached := make([]models.Incident, limit)

	cache.EXPECT().GetTop(gomock.Any()).Return(cached, nil)
	repo.EXPECT().GetTop(gomock.Any(), gomock.Any()).Times(0)
	cache.EXPECT().SetTop(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	result, err := service.GetIncidents(context.Background(), page, limit)

	assert.NoError(t, err)
	assert.Equal(t, cached, result)
}

func TestTop_FromRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()

	service := services.NewIncidentService(repo, cache, logger)

	page := 1
	limit := 10
	rows := make([]models.Incident, 20)

	cache.EXPECT().GetTop(gomock.Any()).Return([]models.Incident{}, nil)
	repo.EXPECT().GetTop(gomock.Any(), topLimit).Return(rows, nil)
	cache.EXPECT().SetTop(gomock.Any(), rows, time.Minute)

	res, err := service.GetIncidents(context.Background(), page, limit)

	assert.NoError(t, err)
	assert.Equal(t, rows, res)
}

func TestGetIncidents_Top_CacheError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()

	service := services.NewIncidentService(repo, cache, logger)

	cache.EXPECT().GetTop(gomock.Any()).Return(nil, errors.New("redis down")) //имитируем ошибку редиса

	res, err := service.GetIncidents(context.Background(), page, limit)

	assert.Error(t, err)
	assert.Nil(t, res)
}

// набор тестов для GetIncidentByID
func TestGetIncidentByID_FromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()

	services := services.NewIncidentService(repo, cache, logger)
	val := &models.Incident{}
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(val, nil)

	result, err := services.GetIncidentByID(t.Context(), 0)
	assert.NoError(t, err)
	assert.Equal(t, &models.Incident{}, result)

}

func TestGetIncidentByID_FromRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()

	services := services.NewIncidentService(repo, cache, logger)
	val := &models.Incident{}
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
	repo.EXPECT().GetIncidentByID(gomock.Any(), 0).Return(val, nil)
	cache.EXPECT().Set(gomock.Any(), fmt.Sprintf("incident:%d", id), val, deafaultTTL*time.Minute).Return(nil)

	result, err := services.GetIncidentByID(t.Context(), 0)
	assert.NoError(t, err)
	assert.Equal(t, val, result)
}

// набор тестов для GetIncidentStat
func TestGetIncidentStat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	fromTime := time.Now().Add(-time.Duration(StatsTimeWindowMinutes) * time.Minute)
	services := services.NewIncidentService(repo, cache, logger)
	users := &models.UsersInDangerousZones{}
	repo.EXPECT().GetIncidentStat(gomock.Any(), &fromTime).Return(users, nil)

	res, err := services.GetIncidentStat(t.Context(), StatsTimeWindowMinutes)
	assert.NoError(t, err)
	assert.Equal(t, users, res)
}

func TestGetIncidentStat_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	fromTime := time.Now().Add(-time.Duration(StatsTimeWindowMinutes) * time.Minute)
	services := services.NewIncidentService(repo, cache, logger)
	users := &models.UsersInDangerousZones{}
	repo.EXPECT().GetIncidentStat(gomock.Any(), &fromTime).Return(users, errors.New("db is down")) //имитируем ошибку базы данных

	res, err := services.GetIncidentStat(t.Context(), StatsTimeWindowMinutes)
	assert.Error(t, err)
	assert.Nil(t, res)
}

// набор тестов для CreateIncidents
func TestCreateIncidents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)
	dto := &dto.IncidentDTO{}
	repo.EXPECT().CreateIncidents(gomock.Any(), dto).Return(nil)

	err := services.CreateIncidents(t.Context(), dto)
	assert.NoError(t, err)
}

func TestCreateIncidents_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)
	dto := &dto.IncidentDTO{}
	repo.EXPECT().CreateIncidents(gomock.Any(), dto).Return(errors.New("db is down")) //имитируем ошибку базы данных

	err := services.CreateIncidents(t.Context(), dto)
	assert.Error(t, err)
}

// набор тестов для UpdateZones
func TestUpdateZones(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)
	dto := &dto.IncidentDTO{}
	repo.EXPECT().UpdateZones(gomock.Any(), dto).Return(nil)

	err := services.UpdateZones(t.Context(), dto)
	assert.NoError(t, err)
}
func TestUpdateZones_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)
	dto := &dto.IncidentDTO{}
	repo.EXPECT().UpdateZones(gomock.Any(), dto).Return(errors.New("db is down")) //имитируем ошибку базы данных

	err := services.UpdateZones(t.Context(), dto)
	assert.Error(t, err)
}

// набор тестов для UpdateIncidentsByID
func TestUpdateIncidentsByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)
	dto := &dto.IncidentDTO{}
	repo.EXPECT().UpdateIncidentsByID(gomock.Any(), dto, id).Return(nil)
	cache.EXPECT().Del(gomock.Any(), id).Return(nil)

	err := services.UpdateIncidentsByID(t.Context(), dto, id)
	assert.NoError(t, err)
}

func TestUpdateIncidentsByID_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)
	dto := &dto.IncidentDTO{}
	repo.EXPECT().UpdateIncidentsByID(gomock.Any(), dto, id).Return(errors.New("db is down")) //имитируем ошибку базы данных

	err := services.UpdateIncidentsByID(t.Context(), dto, id)
	assert.Error(t, err)
}
func TestUpdateIncidentsByID_CacheError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)
	dto := &dto.IncidentDTO{}
	repo.EXPECT().UpdateIncidentsByID(gomock.Any(), dto, id).Return(nil)
	cache.EXPECT().Del(gomock.Any(), id).Return(errors.New("rdb is down")) //имитируем ошибку редиса
	err := services.UpdateIncidentsByID(t.Context(), dto, id)
	assert.Error(t, err)
}

// набор тестов для DeleteIncidentsByID
func TestDeleteIncidentsByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)

	repo.EXPECT().DeleteIncidentsByID(gomock.Any(), id).Return(nil)
	cache.EXPECT().Del(gomock.Any(), id).Return(nil)

	err := services.DeleteIncidentsByID(t.Context(), id)
	assert.NoError(t, err)
}

func TestDeleteIncidentsByID_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)

	repo.EXPECT().DeleteIncidentsByID(gomock.Any(), id).Return(errors.New("db is down")) //имитируем ошибку базы данных

	err := services.DeleteIncidentsByID(t.Context(), id)
	assert.Error(t, err)
}
func TestDeleteIncidentsByID_CacheError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_services.NewMockincidentRepository(ctrl)
	cache := mock_services.NewMockincidentCache(ctrl)
	logger := zap.NewNop()
	services := services.NewIncidentService(repo, cache, logger)

	repo.EXPECT().DeleteIncidentsByID(gomock.Any(), id).Return(nil)
	cache.EXPECT().Del(gomock.Any(), id).Return(errors.New("rdb is down")) //имитируем ошибку редиса

	err := services.DeleteIncidentsByID(t.Context(), id)
	assert.Error(t, err)
}
