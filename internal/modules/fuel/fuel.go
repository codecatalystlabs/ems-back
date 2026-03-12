package fuel

import (
	fuelapp "dispatch/internal/modules/fuel/application"
	"dispatch/internal/modules/fuel/infrastructure"
	"dispatch/internal/modules/fuel/infrastructure/http"
	rbacapp "dispatch/internal/modules/rbac/application"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps, rbacSvc *rbacapp.Service) {
	repo := infrastructure.NewRepository(deps.DB)
	service := fuelapp.NewService(repo, deps.Logger)
	handler := http.NewHandler(service)
	group := deps.Router.Group("/fuel")
	http.RegisterRoutes(group, handler, rbacSvc)
}
