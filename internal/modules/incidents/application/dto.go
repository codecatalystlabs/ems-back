package application

type CreateIncidentRequest struct {
	SourceChannel   string   `json:"source_channel" binding:"required" enums:"SMS,USSD,CALL,MOBILE_APP,WEB_PORTAL,FACILITY_REFERRAL"`
	CallerName      *string  `json:"caller_name,omitempty"`
	CallerPhone     *string  `json:"caller_phone,omitempty"`
	PatientName     *string  `json:"patient_name,omitempty"`
	PatientPhone    *string  `json:"patient_phone,omitempty"`
	PatientAgeGroup *string  `json:"patient_age_group,omitempty" enums:"INFANT,CHILD,ADOLESCENT,ADULT,ELDERLY"`
	PatientSex      *string  `json:"patient_sex,omitempty" enums:"MALE,FEMALE,OTHER,UNKNOWN"`
	// IncidentTypeID can be a UUID, a code, or a name; it will be resolved server-side.
	IncidentTypeID  string   `json:"incident_type_id" binding:"required"`
	Summary         *string  `json:"summary,omitempty"`
	Description     *string  `json:"description,omitempty"`
	// DistrictID can be a UUID, code, or name; it will be resolved server-side.
	DistrictID *string `json:"district_id,omitempty"`
	// FacilityID can be a UUID or facility code (UID); it will be resolved server-side.
	FacilityID *string  `json:"facility_id,omitempty"`
	Village    *string  `json:"village,omitempty"`
	Parish     *string  `json:"parish,omitempty"`
	Subcounty  *string  `json:"subcounty,omitempty"`
	Landmark   *string  `json:"landmark,omitempty"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
}

type UpdateIncidentRequest struct {
	CallerName      *string `json:"caller_name,omitempty"`
	CallerPhone     *string `json:"caller_phone,omitempty"`
	PatientName     *string `json:"patient_name,omitempty"`
	PatientPhone    *string `json:"patient_phone,omitempty"`
	PatientAgeGroup *string `json:"patient_age_group,omitempty"`
	PatientSex      *string `json:"patient_sex,omitempty"`
	IncidentTypeID  *string `json:"incident_type_id,omitempty"`
	SeverityLevelID *string `json:"severity_level_id,omitempty"`
	PriorityLevelID *string `json:"priority_level_id,omitempty"`
	Summary         *string `json:"summary,omitempty"`
	Description     *string `json:"description,omitempty"`
	DistrictID      *string `json:"district_id,omitempty"`
	FacilityID      *string `json:"facility_id,omitempty"`
	Village         *string `json:"village,omitempty"`
	Parish          *string `json:"parish,omitempty"`
	Subcounty       *string `json:"subcounty,omitempty"`
	Landmark        *string `json:"landmark,omitempty"`
	Status          *string `json:"status,omitempty" enums:"NEW,PENDING_VERIFICATION,VERIFIED,AWAITING_ASSIGNMENT,ASSIGNED,ENROUTE,AT_SCENE,TRANSPORTING,COMPLETED,CANCELLED,ESCALATED,REJECTED"`
	Verification    *string `json:"verification_status,omitempty" enums:"PENDING,VERIFIED,REJECTED"`
}
