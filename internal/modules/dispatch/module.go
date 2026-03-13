package dispatch

import (
	dispatchapp "dispatch/internal/modules/dispatch/application"
	"dispatch/internal/modules/dispatch/infrastructure"
	"dispatch/internal/modules/dispatch/infrastructure/http"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	service := dispatchapp.NewService(repo, deps.Bus)
	h := http.NewHandler(service)
	group := deps.Router.Group("/dispatch")
	http.RegisterRoutes(group, h)
}
