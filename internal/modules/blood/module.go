package blood

import (
	bloodapp "dispatch/internal/modules/blood/application"
	"dispatch/internal/modules/blood/infrastructure"
	"dispatch/internal/modules/blood/infrastructure/http"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	service := bloodapp.NewService(repo, deps.Bus, deps.Logger)
	handler := http.NewHandler(service)
	group := deps.Router.Group("/blood")
	http.RegisterRoutes(group, handler)
}
