package dto

import "time"

type CreateBloodRequisitionRequest struct {
	IncidentID            *string    `json:"incident_id"`
	RequestingFacilityID  *string    `json:"requesting_facility_id"`
	PatientName           string     `json:"patient_name"`
	PatientIdentifier     string     `json:"patient_identifier"`
	ClinicalSummary       string     `json:"clinical_summary" binding:"required"`
	Diagnosis             string     `json:"diagnosis"`
	Indication            string     `json:"indication"`
	ParitySummary         string     `json:"parity_summary"`
	BloodGroupCode        string     `json:"blood_group_code" binding:"required"`
	BloodProductCode      string     `json:"blood_product_code" binding:"required"`
	UnitsRequested        int        `json:"units_requested" binding:"required,min=1"`
	UrgencyLevel          string     `json:"urgency_level"`
	ReporterPhone         string     `json:"reporter_phone"`
	DestinationFacilityID *string    `json:"destination_facility_id"`
	DestinationLat        *float64   `json:"destination_lat"`
	DestinationLon        *float64   `json:"destination_lon"`
	RequestedByUserID     *string    `json:"requested_by_user_id"`
	ExpiresAt             *time.Time `json:"expires_at"`
}

type CreateBloodOfferRequest struct {
	BloodRequisitionID string     `json:"blood_requisition_id" binding:"required"`
	InventorySiteID    string     `json:"inventory_site_id" binding:"required"`
	BloodGroupCode     string     `json:"blood_group_code" binding:"required"`
	BloodProductCode   string     `json:"blood_product_code" binding:"required"`
	UnitsOffered       int        `json:"units_offered" binding:"required,min=1"`
	ReservedUntil      *time.Time `json:"reserved_until"`
	ContactPersonName  string     `json:"contact_person_name"`
	ContactPhone       string     `json:"contact_phone"`
	Notes              string     `json:"notes"`
	OfferedByUserID    *string    `json:"offered_by_user_id"`
}

type AssignBloodPickupRequest struct {
	BloodRequisitionID      string  `json:"blood_requisition_id" binding:"required"`
	BloodRequisitionOfferID *string `json:"blood_requisition_offer_id"`
	VehicleType             string  `json:"vehicle_type" binding:"required"`
	AmbulanceID             *string `json:"ambulance_id"`
	DispatchAssignmentID    *string `json:"dispatch_assignment_id"`
	AssignedDriverUserID    *string `json:"assigned_driver_user_id"`
	AssignedByUserID        *string `json:"assigned_by_user_id"`
	PickupSiteID            *string `json:"pickup_site_id"`
	DestinationFacilityID   *string `json:"destination_facility_id"`
	Notes                   string  `json:"notes"`
}
