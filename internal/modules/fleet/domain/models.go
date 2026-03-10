package domain

import "time"

type Ambulance struct {
	ID                 string    `json:"id"`
	Code               *string   `json:"code,omitempty"`
	PlateNumber        string    `json:"plate_number"`
	VIN                *string   `json:"vin,omitempty"`
	Make               *string   `json:"make,omitempty"`
	Model              *string   `json:"model,omitempty"`
	YearOfManufacture  *int      `json:"year_of_manufacture,omitempty"`
	CategoryID         string    `json:"category_id"`
	OwnershipType      *string   `json:"ownership_type,omitempty"`
	StationFacilityID  *string   `json:"station_facility_id,omitempty"`
	DistrictID         *string   `json:"district_id,omitempty"`
	Status             string    `json:"status"`
	DispatchReadiness  string    `json:"dispatch_readiness"`
	GPSLat             *float64  `json:"gps_lat,omitempty"`
	GPSLon             *float64  `json:"gps_lon,omitempty"`
	LastSeenAt         *time.Time `json:"last_seen_at,omitempty"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

