package incidents

import (
	incapp "dispatch/internal/modules/incidents/application"
	"dispatch/internal/modules/incidents/infrastructure"
	"dispatch/internal/modules/incidents/infrastructure/http"
	"dispatch/internal/shared/types"

	rbacapp "dispatch/internal/modules/rbac/application"
)

func Register(deps types.ModuleDeps, rbacSvc *rbacapp.Service) {
	repo := infrastructure.NewRepository(deps.DB)
	service := incapp.NewService(repo, deps.Bus, deps.Logger)
	handler := http.NewHandler(service)
	group := deps.Router.Group("/incidents")
	http.RegisterRoutes(group, handler, rbacSvc)
}
