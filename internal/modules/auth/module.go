package auth

import (
	authapp "dispatch/internal/modules/auth/application"
	"dispatch/internal/modules/auth/infrastructure"
	"dispatch/internal/modules/auth/infrastructure/http"
	middleware "dispatch/internal/modules/auth/middleware"
	platformauth "dispatch/internal/platform/auth"
	"dispatch/internal/shared/types"
)

func Register(deps types.ModuleDeps) {
	repo := infrastructure.NewRepository(deps.DB)
	jwt := platformauth.NewJWTManager(deps.Config.JWT.Secret, deps.Config.JWT.Issuer, deps.Config.JWT.AccessTTL, deps.Config.JWT.RefreshTTL)
	service := authapp.NewService(repo, jwt, deps.Redis, deps.Logger, deps.Config.JWT.AccessTTL, deps.Config.JWT.RefreshTTL)
	h := http.NewHandler(service)
	group := deps.Router.Group("/auth")
	http.RegisterRoutes(group, h, middleware.AuthMiddleware(deps.Config.JWT.Secret))
}
