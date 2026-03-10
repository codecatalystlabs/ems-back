package application

import (
	"context"

	rbacdomain "dispatch/internal/modules/rbac/domain"
)

type Repository interface {
	ListPermissionGrants(ctx context.Context, userID string) ([]rbacdomain.PermissionGrant, error)
}
