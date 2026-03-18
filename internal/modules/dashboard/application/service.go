package application

import (
	"context"
	"time"

	dashboarddomain "dispatch/internal/modules/dashboard/domain"
)

type Service struct {
	repo    Repository
	cache   *CacheService
	refresh *RefreshService
}

func NewService(repo Repository, cache *CacheService, refresh *RefreshService) *Service {
	return &Service{
		repo:    repo,
		cache:   cache,
		refresh: refresh,
	}
}

func (s *Service) GetDashboard(ctx context.Context, q DashboardQuery) (dashboarddomain.DashboardResponse, error) {
	key := dashboardKey(q.DateFrom, q.DateTo, q.DistrictID, q.FacilityID)

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, key)
		if err == nil && cached != nil {
			return *cached, nil
		}
	}

	if s.refresh != nil {
		_ = s.refresh.RefreshAll(ctx)
	}

	filters := dashboarddomain.DashboardFilters{}
	if q.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", q.DateFrom); err == nil {
			filters.DateFrom = &t
		}
	}
	if q.DateTo != "" {
		if t, err := time.Parse("2006-01-02", q.DateTo); err == nil {
			end := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filters.DateTo = &end
		}
	}
	if q.DistrictID != "" {
		filters.DistrictID = &q.DistrictID
	}
	if q.FacilityID != "" {
		filters.FacilityID = &q.FacilityID
	}

	out, err := s.repo.GetDashboard(ctx, filters)
	if err != nil {
		return dashboarddomain.DashboardResponse{}, err
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, key, out)
	}

	return out, nil
}
