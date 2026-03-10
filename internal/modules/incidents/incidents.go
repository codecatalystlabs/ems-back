package incidents

import (
	incapp "dispatch/internal/modules/incidents/application"
	"dispatch/internal/modules/incidents/infrastructure"
	"dispatch/internal/modules/incidents/infrastructure/http"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	service := incapp.NewService(repo, deps.Logger)
	handler := http.NewHandler(service)
	group := deps.Router.Group("/incidents")
	http.RegisterRoutes(group, handler)
}

