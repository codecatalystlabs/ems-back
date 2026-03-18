package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	dashboarddomain "dispatch/internal/modules/dashboard/domain"
)

type CacheService struct {
	redis *redis.Client
	ttl   time.Duration
}

func NewCacheService(redis *redis.Client) *CacheService {
	return &CacheService{
		redis: redis,
		ttl:   2 * time.Minute,
	}
}

func dashboardKey(dateFrom, dateTo, districtID, facilityID string) string {
	return fmt.Sprintf(
		"dashboard:ems_overview:v1:from=%s:to=%s:district=%s:facility=%s",
		dateFrom, dateTo, districtID, facilityID,
	)
}

func (s *CacheService) Get(ctx context.Context, key string) (*dashboarddomain.DashboardResponse, error) {
	if s.redis == nil {
		return nil, nil
	}

	val, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var out dashboarddomain.DashboardResponse
	if err := json.Unmarshal([]byte(val), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *CacheService) Set(ctx context.Context, key string, value dashboarddomain.DashboardResponse) error {
	if s.redis == nil {
		return nil
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, b, s.ttl).Err()
}

func (s *CacheService) InvalidateAllOverview(ctx context.Context) error {
	if s.redis == nil {
		return nil
	}
	keys, err := s.redis.Keys(ctx, "dashboard:ems_overview:v1:*").Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return s.redis.Del(ctx, keys...).Err()
}
