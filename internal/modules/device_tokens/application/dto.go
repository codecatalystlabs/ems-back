package application

import platformdb "dispatch/internal/platform/db"

type RegisterDeviceTokenRequest struct {
	UserID    string `json:"user_id" binding:"required,uuid"`
	DeviceID  string `json:"device_id"`
	Platform  string `json:"platform" binding:"required"`
	PushToken string `json:"push_token" binding:"required"`
}

type UpdateDeviceTokenRequest struct {
	DeviceID  *string `json:"device_id"`
	Platform  *string `json:"platform"`
	PushToken *string `json:"push_token"`
	IsActive  *bool   `json:"is_active"`
}

type ListDeviceTokensParams struct {
	UserID     *string               `json:"user_id,omitempty"`
	Platform   *string               `json:"platform,omitempty"`
	IsActive   *bool                 `json:"is_active,omitempty"`
	Pagination platformdb.Pagination `json:"pagination"`
}
