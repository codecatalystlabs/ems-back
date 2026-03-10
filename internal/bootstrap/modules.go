package bootstrap

import (
	// "dispatch/internal/modules/auth"
	// availabilitymod "dispatch/internal/modules/availability"
	// dispatchmod "dispatch/internal/modules/dispatch"
	// fleetmod "dispatch/internal/modules/fleet"
	// incidentmod "dispatch/internal/modules/incidents"
	// rbacmod "dispatch/internal/modules/rbac"
	authmod "dispatch/internal/modules/auth"

	bloodmod "dispatch/internal/modules/blood"
	usermod "dispatch/internal/modules/users"
	"dispatch/internal/shared/types"
)

func RegisterModules(deps types.ModuleDeps) {
	authmod.Register(deps)
	usermod.Register(deps)
	// rbacmod.Register(deps)
	bloodmod.Register(deps)
	// fleetmod.Register(deps)
	// availabilitymod.Register(deps)
	// incidentmod.Register(deps)
	// dispatchmod.Register(deps)
}
