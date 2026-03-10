package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	rbacdomain "dispatch/internal/modules/rbac/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPermissionGrants(ctx context.Context, userID string) ([]rbacdomain.PermissionGrant, error) {
	rows, err := r.db.Query(ctx, `
		SELECT ur.user_id, r.code, p.code, ur.scope_type, ur.scope_id::text
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		JOIN role_permissions rp ON rp.role_id = r.id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ur.user_id = $1 AND ur.active = TRUE
		ORDER BY r.code, p.code
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grants []rbacdomain.PermissionGrant
	for rows.Next() {
		var g rbacdomain.PermissionGrant
		var scopeID *string
		if err := rows.Scan(&g.UserID, &g.RoleCode, &g.PermCode, &g.ScopeType, &scopeID); err != nil {
			return nil, err
		}
		g.ScopeID = scopeID
		grants = append(grants, g)
	}
	return grants, rows.Err()
}
