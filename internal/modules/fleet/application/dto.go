package application

type CreateAmbulanceRequest struct {
	Code              *string `json:"code,omitempty"`
	PlateNumber       string  `json:"plate_number" binding:"required"`
	VIN               *string `json:"vin,omitempty"`
	Make              *string `json:"make,omitempty"`
	Model             *string `json:"model,omitempty"`
	YearOfManufacture *int    `json:"year_of_manufacture,omitempty"`
	CategoryID        string  `json:"category_id" binding:"required"`
	OwnershipType     *string `json:"ownership_type,omitempty"`
	StationFacilityID *string `json:"station_facility_id,omitempty"`
	DistrictID        *string `json:"district_id,omitempty"`
	Status            *string `json:"status,omitempty" enums:"AVAILABLE,RESERVED,ASSIGNED,ENROUTE,AT_SCENE,TRANSPORTING,RETURNING,MAINTENANCE,BREAKDOWN,OFFLINE,RETIRED"`
	DispatchReadiness *string `json:"dispatch_readiness,omitempty" enums:"DISPATCHABLE,RESTRICTED,NOT_DISPATCHABLE"`
}

type UpdateAmbulanceRequest struct {
	Code              *string `json:"code,omitempty"`
	VIN               *string `json:"vin,omitempty"`
	Make              *string `json:"make,omitempty"`
	Model             *string `json:"model,omitempty"`
	YearOfManufacture *int    `json:"year_of_manufacture,omitempty"`
	CategoryID        *string `json:"category_id,omitempty"`
	OwnershipType     *string `json:"ownership_type,omitempty"`
	StationFacilityID *string `json:"station_facility_id,omitempty"`
	DistrictID        *string `json:"district_id,omitempty"`
	Status            *string `json:"status,omitempty" enums:"AVAILABLE,RESERVED,ASSIGNED,ENROUTE,AT_SCENE,TRANSPORTING,RETURNING,MAINTENANCE,BREAKDOWN,OFFLINE,RETIRED"`
	DispatchReadiness *string `json:"dispatch_readiness,omitempty" enums:"DISPATCHABLE,RESTRICTED,NOT_DISPATCHABLE"`
}
