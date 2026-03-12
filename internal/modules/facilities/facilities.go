package facilities

import (
	facapp "dispatch/internal/modules/facilities/application"
	"dispatch/internal/modules/facilities/infrastructure"
	"dispatch/internal/modules/facilities/infrastructure/http"
	rbacapp "dispatch/internal/modules/rbac/application"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps, rbacSvc *rbacapp.Service) {
	repo := infrastructure.NewRepository(deps.DB)
	service := facapp.NewService(repo, deps.Logger)
	handler := http.NewHandler(service)
	group := deps.Router.Group("/facilities")
	http.RegisterRoutes(group, handler, rbacSvc)
}
