package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`

	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	Platform   string `json:"platform"`
	PushToken  string `json:"push_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
	LogoutAll    bool   `json:"logout_all"`
}

type AuthResponse struct {
	AccessToken           string           `json:"access_token"`
	RefreshToken          string           `json:"refresh_token"`
	AccessTokenExpiresIn  int64            `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int64            `json:"refresh_token_expires_in"`
	Session               SessionResponse  `json:"session"`
	User                  AuthUserResponse `json:"user"`
	Permissions           []string         `json:"permissions"`
}

type AuthUserResponse struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Status   string   `json:"status"`
	Roles    []string `json:"roles"`
}

type SessionResponse struct {
	ID             string `json:"id"`
	DeviceID       string `json:"device_id"`
	DeviceName     string `json:"device_name"`
	IPAddress      string `json:"ip_address"`
	UserAgent      string `json:"user_agent"`
	LastActivityAt string `json:"last_activity_at"`
	ExpiresAt      string `json:"expires_at"`
}

type PermissionGrantDTO struct {
	UserID    string  `json:"user_id"`
	RoleCode  string  `json:"role_code"`
	PermCode  string  `json:"perm_code"`
	ScopeType string  `json:"scope_type"`
	ScopeID   *string `json:"scope_id,omitempty"`
}
