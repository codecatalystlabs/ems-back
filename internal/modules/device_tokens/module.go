package devicetokens

import (
	deviceapp "dispatch/internal/modules/device_tokens/application"
	"dispatch/internal/modules/device_tokens/infrastructure"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	service := deviceapp.NewService(repo)
	handler := infrastructure.NewHandler(service)

	group := deps.Router.Group("/device-tokens")
	infrastructure.RegisterRoutes(group, handler)
}
