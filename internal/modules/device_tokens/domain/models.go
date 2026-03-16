package domain

import "time"

type DeviceToken struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	DeviceID   string     `json:"device_id"`
	Platform   string     `json:"platform"`
	PushToken  string     `json:"push_token"`
	IsActive   bool       `json:"is_active"`
	LastSeenAt *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
