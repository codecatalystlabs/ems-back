package dashboard

import (
	dashboardapp "dispatch/internal/modules/dashboard/application"
	"dispatch/internal/modules/dashboard/infrastructure"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	cache := dashboardapp.NewCacheService(deps.Redis)
	refresh := dashboardapp.NewRefreshService(deps.DB)
	service := dashboardapp.NewService(repo, cache, refresh)
	handler := infrastructure.NewHandler(service)

	group := deps.Router.Group("/dashboard")
	infrastructure.RegisterRoutes(group, handler)
}
