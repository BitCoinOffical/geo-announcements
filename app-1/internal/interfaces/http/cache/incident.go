package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
	"github.com/redis/go-redis/v9"
)

const topKey = "incidents:top"

type IncidentCache struct {
	rdb *redis.Client
}

func NewIncidentCache(rdb *redis.Client) *IncidentCache {
	return &IncidentCache{rdb: rdb}
}

func (c *IncidentCache) GetTop(ctx context.Context) ([]models.Incident, error) {
	res, err := c.rdb.Get(ctx, topKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var model []models.Incident
	if err := json.Unmarshal([]byte(res), &model); err != nil {
		return nil, err
	}
	return model, nil
}

func (c *IncidentCache) GetAll(ctx context.Context, page, limit int) ([]models.Incident, error) {
	key := fmt.Sprintf("incident-page%d:-limit:%d", page, limit)
	res, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var model []models.Incident
	if err := json.Unmarshal([]byte(res), &model); err != nil {
		return nil, err
	}
	return model, nil
}

func (c *IncidentCache) Get(ctx context.Context, id int) (*models.Incident, error) {
	key := fmt.Sprintf("incident:%d", id)
	res, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var model models.Incident
	if err := json.Unmarshal([]byte(res), &model); err != nil {
		return nil, err
	}
	return &model, nil
}

func (c *IncidentCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

func (c *IncidentCache) SetTop(ctx context.Context, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, topKey, data, ttl).Err()
}

func (c *IncidentCache) Del(ctx context.Context, id int) error {
	key := fmt.Sprintf("incident:%d", id)
	return c.rdb.Del(ctx, key).Err()
}
