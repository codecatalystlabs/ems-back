package application

import (
	"context"

	dashboarddomain "dispatch/internal/modules/dashboard/domain"
)

type Repository interface {
	GetDashboard(ctx context.Context, filters dashboarddomain.DashboardFilters) (dashboarddomain.DashboardResponse, error)
}
