package auth

import (
	authapp "dispatch/internal/modules/auth/application"
	authinfra "dispatch/internal/modules/auth/infrastructure"
	authhttp "dispatch/internal/modules/auth/infrastructure/http"
	middleware "dispatch/internal/modules/auth/middleware"
	rbacinfra "dispatch/internal/modules/rbac/infrastructure"
	platformauth "dispatch/internal/platform/auth"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	authRepo := authinfra.NewRepository(deps.DB)
	rbacRepo := rbacinfra.NewRepository(deps.DB)

	jwt := platformauth.NewJWTManager(
		deps.Config.JWT.Secret,
		deps.Config.JWT.Issuer,
		deps.Config.JWT.AccessTTL,
		deps.Config.JWT.RefreshTTL,
	)

	service := authapp.NewService(
		authRepo,
		rbacRepo,
		jwt,
		deps.Redis,
		deps.Logger,
		deps.Config.JWT.AccessTTL,
		deps.Config.JWT.RefreshTTL,
	)

	h := authhttp.NewHandler(service)
	group := deps.Router.Group("/auth")
	authhttp.RegisterRoutes(group, h, middleware.AuthMiddleware(deps.Config.JWT.Secret))
}
