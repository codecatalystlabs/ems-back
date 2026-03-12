package application

import "time"

type CreateFuelLogRequest struct {
	AmbulanceID string     `json:"ambulance_id" binding:"required,uuid"`
	FuelType    *string    `json:"fuel_type,omitempty"`
	Liters      float64    `json:"liters" binding:"required,gt=0"`
	Cost        *float64   `json:"cost,omitempty"`
	OdometerKM  *int       `json:"odometer_km,omitempty"`
	StationName *string    `json:"station_name,omitempty"`
	FilledAt    *time.Time `json:"filled_at,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
}

type UpdateFuelLogRequest struct {
	FuelType    *string    `json:"fuel_type,omitempty"`
	Liters      *float64   `json:"liters,omitempty"`
	Cost        *float64   `json:"cost,omitempty"`
	OdometerKM  *int       `json:"odometer_km,omitempty"`
	StationName *string    `json:"station_name,omitempty"`
	FilledAt    *time.Time `json:"filled_at,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
}
