package dto

import platformdb "dispatch/internal/platform/db"

type ListDistrictsParams struct {
	Pagination platformdb.Pagination `json:"pagination"`
}

type ListSubcountiesParams struct {
	DistrictID *string               `json:"district_id,omitempty"`
	Pagination platformdb.Pagination `json:"pagination"`
}

type ListFacilitiesParams struct {
	DistrictID  *string               `json:"district_id,omitempty"`
	SubcountyID *string               `json:"subcounty_id,omitempty"`
	LevelID     *string               `json:"level_id,omitempty"`
	Pagination  platformdb.Pagination `json:"pagination"`
}
