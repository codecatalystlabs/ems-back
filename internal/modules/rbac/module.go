package rbac

import (
	rbacapp "dispatch/internal/modules/rbac/application"
	"dispatch/internal/modules/rbac/infrastructure"
	"dispatch/internal/modules/rbac/infrastructure/http"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) *rbacapp.Service {
	repo := infrastructure.NewRepository(deps.DB)
	service := rbacapp.NewService(repo, deps.Redis, deps.Logger)
	h := http.NewHandler(service)
	group := deps.Router.Group("/rbac")
	http.RegisterRoutes(group, h)
	return service
}
