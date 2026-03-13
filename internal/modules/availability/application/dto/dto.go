package dto

import (
	platformdb "dispatch/internal/platform/db"
)

type CreateShiftRequest struct {
	UserID     string  `json:"user_id" binding:"required,uuid"`
	ShiftDate  string  `json:"shift_date" binding:"required"`
	StartsAt   string  `json:"starts_at" binding:"required"`
	EndsAt     string  `json:"ends_at" binding:"required"`
	ShiftType  string  `json:"shift_type"`
	DistrictID *string `json:"district_id"`
	FacilityID *string `json:"facility_id"`
	Status     string  `json:"status"`
	CreatedBy  *string `json:"created_by"`
}

type UpdateShiftRequest struct {
	ShiftDate  *string `json:"shift_date"`
	StartsAt   *string `json:"starts_at"`
	EndsAt     *string `json:"ends_at"`
	ShiftType  *string `json:"shift_type"`
	DistrictID *string `json:"district_id"`
	FacilityID *string `json:"facility_id"`
	Status     *string `json:"status"`
}

type UpsertAvailabilityRequest struct {
	UserID                      string  `json:"user_id" binding:"required,uuid"`
	AvailabilityStatus          string  `json:"availability_status" binding:"required"`
	Dispatchable                bool    `json:"dispatchable"`
	CurrentIncidentID           *string `json:"current_incident_id"`
	CurrentDispatchAssignmentID *string `json:"current_dispatch_assignment_id"`
	CurrentAmbulanceID          *string `json:"current_ambulance_id"`
	LastSeenAt                  *string `json:"last_seen_at"`
	Source                      string  `json:"source"`
	Notes                       string  `json:"notes"`
	UpdatedBy                   *string `json:"updated_by"`
}

type CreatePresenceLogRequest struct {
	UserID    string   `json:"user_id" binding:"required,uuid"`
	Channel   string   `json:"channel" binding:"required"`
	IPAddress *string  `json:"ip_address"`
	UserAgent string   `json:"user_agent"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type ListShiftsParams struct {
	UserID     *string               `json:"user_id,omitempty"`
	DistrictID *string               `json:"district_id,omitempty"`
	FacilityID *string               `json:"facility_id,omitempty"`
	ShiftDate  *string               `json:"shift_date,omitempty"`
	Status     *string               `json:"status,omitempty"`
	Pagination platformdb.Pagination `json:"pagination"`
}

type ListAvailabilityParams struct {
	Status       *string               `json:"status,omitempty"`
	Dispatchable *bool                 `json:"dispatchable,omitempty"`
	Pagination   platformdb.Pagination `json:"pagination"`
}

type ListPresenceParams struct {
	UserID     *string               `json:"user_id,omitempty"`
	Channel    *string               `json:"channel,omitempty"`
	Pagination platformdb.Pagination `json:"pagination"`
}
