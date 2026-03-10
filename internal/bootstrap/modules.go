package bootstrap

import (
	// availabilitymod "dispatch/internal/modules/availability"
	// dispatchmod "dispatch/internal/modules/dispatch"
	// fleetmod "dispatch/internal/modules/fleet"
	// incidentmod "dispatch/internal/modules/incidents"
	authmod "dispatch/internal/modules/auth"
	rbacmod "dispatch/internal/modules/rbac"

	bloodmod "dispatch/internal/modules/blood"
	usermod "dispatch/internal/modules/users"
	"dispatch/internal/shared/types"

	authmiddleware "dispatch/internal/modules/auth/middleware"
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"
)

func RegisterModules(deps types.ModuleDeps) {
	authmod.Register(deps)
	rbacSvc := rbacmod.Register(deps)

	secured := deps.Router.Group("")
	secured.Use(authmiddleware.AuthMiddleware(deps.Config.JWT.Secret), rbacmiddleware.ScopeContextMiddleware())

	usersGroup := secured.Group("/users")
	usersGroup.Use(rbacmiddleware.RequirePermission(rbacSvc, "users.read"))
	_ = usersGroup

	fleetGroup := secured.Group("/ambulances")
	fleetGroup.Use(rbacmiddleware.RequirePermission(rbacSvc, "fleet.read"))
	_ = fleetGroup

	usermod.Register(types.ModuleDeps{Router: secured, DB: deps.DB, Redis: deps.Redis, Logger: deps.Logger, Bus: deps.Bus, Config: deps.Config})
	// availabilitymod.Register(types.ModuleDeps{Router: secured, DB: deps.DB, Redis: deps.Redis, Logger: deps.Logger, Bus: deps.Bus, Config: deps.Config})
	// incidentmod.Register(types.ModuleDeps{Router: secured, DB: deps.DB, Redis: deps.Redis, Logger: deps.Logger, Bus: deps.Bus, Config: deps.Config})
	// dispatchmod.Register(types.ModuleDeps{Router: secured, DB: deps.DB, Redis: deps.Redis, Logger: deps.Logger, Bus: deps.Bus, Config: deps.Config})
	bloodmod.Register(types.ModuleDeps{Router: secured, DB: deps.DB, Redis: deps.Redis, Logger: deps.Logger, Bus: deps.Bus, Config: deps.Config})
	// fleetmod.Register(types.ModuleDeps{Router: secured, DB: deps.DB, Redis: deps.Redis, Logger: deps.Logger, Bus: deps.Bus, Config: deps.Config})
}
