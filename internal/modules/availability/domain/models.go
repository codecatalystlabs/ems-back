// File: internal/modules/availability/domain/models.go
package domain

import "time"

type UserShift struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ShiftDate  string    `json:"shift_date"`
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
	ShiftType  string    `json:"shift_type"`
	DistrictID *string   `json:"district_id,omitempty"`
	FacilityID *string   `json:"facility_id,omitempty"`
	Status     string    `json:"status"`
	CreatedBy  *string   `json:"created_by,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserAvailability struct {
	ID                          string     `json:"id"`
	UserID                      string     `json:"user_id"`
	AvailabilityStatus          string     `json:"availability_status"`
	Dispatchable                bool       `json:"dispatchable"`
	CurrentIncidentID           *string    `json:"current_incident_id,omitempty"`
	CurrentDispatchAssignmentID *string    `json:"current_dispatch_assignment_id,omitempty"`
	CurrentAmbulanceID          *string    `json:"current_ambulance_id,omitempty"`
	LastSeenAt                  *time.Time `json:"last_seen_at,omitempty"`
	Source                      string     `json:"source"`
	Notes                       string     `json:"notes"`
	UpdatedBy                   *string    `json:"updated_by,omitempty"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}

type UserPresenceLog struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Channel   string    `json:"channel"`
	SeenAt    time.Time `json:"seen_at"`
	IPAddress *string   `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent"`
	Latitude  *float64  `json:"latitude,omitempty"`
	Longitude *float64  `json:"longitude,omitempty"`
}
