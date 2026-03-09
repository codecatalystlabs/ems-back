package auth

import (
	"dispatch/internal/modules/auth/infrastructure/http"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	h := http.NewHandler()
	group := deps.Router.Group("/auth")
	http.RegisterRoutes(group, h)
}
