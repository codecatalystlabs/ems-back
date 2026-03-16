package domain

import "time"

type District struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Region    string    `json:"region"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Subcounty struct {
	ID         string    `json:"id"`
	DistrictID string    `json:"district_id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FacilityLevel struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	RankNo    int       `json:"rank_no"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Facility struct {
	ID                string    `json:"id"`
	Code              string    `json:"code"`
	Name              string    `json:"name"`
	ShortName         string    `json:"short_name"`
	NHFRID            string    `json:"nhfr_id"`
	DistrictID        *string   `json:"district_id,omitempty"`
	DistrictName      string    `json:"district_name,omitempty"`
	SubcountyID       *string   `json:"subcounty_id,omitempty"`
	SubcountyName     string    `json:"subcounty_name,omitempty"`
	LevelID           *string   `json:"level_id,omitempty"`
	LevelName         string    `json:"level_name,omitempty"`
	Ownership         string    `json:"ownership"`
	Phone             string    `json:"phone"`
	Email             string    `json:"email"`
	Address           string    `json:"address"`
	Latitude          *float64  `json:"latitude,omitempty"`
	Longitude         *float64  `json:"longitude,omitempty"`
	IsDispatchStation bool      `json:"is_dispatch_station"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type IncidentType struct {
	ID                string    `json:"id"`
	Code              string    `json:"code"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	RequiresTransport bool      `json:"requires_transport"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
}

type PriorityLevel struct {
	ID                    string    `json:"id"`
	Code                  string    `json:"code"`
	Name                  string    `json:"name"`
	ColorCode             string    `json:"color_code"`
	SortOrder             int       `json:"sort_order"`
	TargetResponseMinutes *int      `json:"target_response_minutes,omitempty"`
	SeverityWeight        int       `json:"severity_weight"`
	EscalationNote        string    `json:"escalation_note"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
}

type SeverityLevel struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type AmbulanceCategory struct {
	ID                   string    `json:"id"`
	Code                 string    `json:"code"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	SupportsMaternal     bool      `json:"supports_maternal"`
	SupportsNeonatal     bool      `json:"supports_neonatal"`
	SupportsTrauma       bool      `json:"supports_trauma"`
	SupportsCriticalCare bool      `json:"supports_critical_care"`
	SupportsReferral     bool      `json:"supports_referral"`
	MinCrewCount         int       `json:"min_crew_count"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
}

type Capability struct {
	ID             string    `json:"id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CapabilityType string    `json:"capability_type"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type TriageQuestion struct {
	ID                string    `json:"id"`
	QuestionnaireID   string    `json:"questionnaire_id"`
	QuestionnaireCode string    `json:"questionnaire_code"`
	Code              string    `json:"code"`
	QuestionText      string    `json:"question_text"`
	ResponseType      string    `json:"response_type"`
	DisplayOrder      int       `json:"display_order"`
	IsRequired        bool      `json:"is_required"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
}
