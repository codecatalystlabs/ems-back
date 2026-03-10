package notifications

import (
	notifapp "dispatch/internal/modules/notifications/application"
	"dispatch/internal/modules/notifications/infrastructure"
	"dispatch/internal/modules/notifications/infrastructure/http"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	service := notifapp.NewService(repo, deps.Logger)
	handler := http.NewHandler(service)
	group := deps.Router.Group("/notifications")
	http.RegisterRoutes(group, handler)
}

