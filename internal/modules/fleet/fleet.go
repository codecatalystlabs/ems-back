package fleet

import (
	fleetapp "dispatch/internal/modules/fleet/application"
	"dispatch/internal/modules/fleet/infrastructure"
	"dispatch/internal/modules/fleet/infrastructure/http"
	rbacapp "dispatch/internal/modules/rbac/application"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps, rbacSvc *rbacapp.Service) {
	repo := infrastructure.NewRepository(deps.DB)
	service := fleetapp.NewService(repo, deps.Logger)
	handler := http.NewHandler(service)
	group := deps.Router.Group("/ambulances")
	http.RegisterRoutes(group, handler, rbacSvc)
}
