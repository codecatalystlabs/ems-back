package reference

import (
	refapp "dispatch/internal/modules/reference/application"
	"dispatch/internal/modules/reference/infrastructure"
	"dispatch/internal/modules/reference/infrastructure/http"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	service := refapp.NewService(repo)
	handler := http.NewHandler(service)

	group := deps.Router.Group("/reference")
	http.RegisterRoutes(group, handler)
}
