package application

import (
	"context"
	"time"

	dashboarddomain "dispatch/internal/modules/dashboard/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDashboard(ctx context.Context, q DashboardQuery) (dashboarddomain.DashboardResponse, error) {
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
	if q.SubcountyID != "" {
		filters.SubcountyID = &q.SubcountyID
	}
	if q.FacilityID != "" {
		filters.FacilityID = &q.FacilityID
	}
	return s.repo.GetDashboard(ctx, filters)
}
