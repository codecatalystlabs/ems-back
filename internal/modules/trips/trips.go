package trips

import (
	tripsapp "dispatch/internal/modules/trips/application"
	"dispatch/internal/modules/trips/infrastructure"
	"dispatch/internal/modules/trips/infrastructure/http"
	rbacapp "dispatch/internal/modules/rbac/application"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps, rbacSvc *rbacapp.Service) {
	repo := infrastructure.NewRepository(deps.DB)
	service := tripsapp.NewService(repo, deps.Logger)
	handler := http.NewHandler(service)
	group := deps.Router.Group("/trips")
	http.RegisterRoutes(group, handler, rbacSvc)
}

