package domain

import "time"

type FuelLog struct {
	ID          string    `json:"id"`
	AmbulanceID string    `json:"ambulance_id"`
	FuelType    *string   `json:"fuel_type,omitempty"`
	Liters      float64   `json:"liters"`
	Cost        *float64  `json:"cost,omitempty"`
	OdometerKM  *int      `json:"odometer_km,omitempty"`
	StationName *string   `json:"station_name,omitempty"`
	FilledAt    time.Time `json:"filled_at"`
	FilledBy    *string   `json:"filled_by,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
