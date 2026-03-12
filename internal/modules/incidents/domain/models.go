package domain

import "time"

type Incident struct {
	ID                 string     `json:"id"`
	IncidentNumber     string     `json:"incident_number"`
	SourceChannel      string     `json:"source_channel"`
	CallerName         *string    `json:"caller_name,omitempty"`
	CallerPhone        *string    `json:"caller_phone,omitempty"`
	PatientName        *string    `json:"patient_name,omitempty"`
	PatientPhone       *string    `json:"patient_phone,omitempty"`
	PatientAgeGroup    *string    `json:"patient_age_group,omitempty"`
	PatientSex         *string    `json:"patient_sex,omitempty"`
	IncidentTypeID     string     `json:"incident_type_id"`
	SeverityLevelID    *string    `json:"severity_level_id,omitempty"`
	PriorityLevelID    *string    `json:"priority_level_id,omitempty"`
	Summary            *string    `json:"summary,omitempty"`
	Description        *string    `json:"description,omitempty"`
	DistrictID         *string    `json:"district_id,omitempty"`
	FacilityID         *string    `json:"facility_id,omitempty"`
	Village            *string    `json:"village,omitempty"`
	Parish             *string    `json:"parish,omitempty"`
	Subcounty          *string    `json:"subcounty,omitempty"`
	Landmark           *string    `json:"landmark,omitempty"`
	Latitude           *float64   `json:"latitude,omitempty"`
	Longitude          *float64   `json:"longitude,omitempty"`
	VerificationStatus string     `json:"verification_status"`
	Status             string     `json:"status"`
	ReportedAt         time.Time  `json:"reported_at"`
	CreatedByUserID    *string    `json:"created_by_user_id,omitempty"`
	TriagedByUserID    *string    `json:"triaged_by_user_id,omitempty"`
	TriagedAt          *time.Time `json:"triaged_at,omitempty"`
	AssignedAt         *time.Time `json:"assigned_at,omitempty"`
	ClosedAt           *time.Time `json:"closed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
