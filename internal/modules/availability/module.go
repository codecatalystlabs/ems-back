// File: internal/modules/availability/module.go
package availability

import (
	availabilityapp "dispatch/internal/modules/availability/application"
	"dispatch/internal/modules/availability/infrastructure"
	"dispatch/internal/modules/availability/infrastructure/http"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	service := availabilityapp.NewService(repo)
	h := http.NewHandler(service)
	group := deps.Router.Group("/availability")
	http.RegisterRoutes(group, h)
}
