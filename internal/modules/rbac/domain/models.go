package domain

type PermissionGrant struct {
	UserID    string  `json:"user_id"`
	RoleCode  string  `json:"role_code"`
	PermCode  string  `json:"perm_code"`
	ScopeType string  `json:"scope_type"`
	ScopeID   *string `json:"scope_id,omitempty"`
}
