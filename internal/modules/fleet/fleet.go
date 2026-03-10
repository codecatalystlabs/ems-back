package fleet

import (
	fleetapp "dispatch/internal/modules/fleet/application"
	"dispatch/internal/modules/fleet/infrastructure"
	"dispatch/internal/modules/fleet/infrastructure/http"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	service := fleetapp.NewService(repo, deps.Logger)
	handler := http.NewHandler(service)
	group := deps.Router.Group("/ambulances")
	http.RegisterRoutes(group, handler)
}

