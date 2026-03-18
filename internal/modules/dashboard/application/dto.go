package application

type DashboardQuery struct {
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
	DistrictID string `form:"district_id"`
	FacilityID string `form:"facility_id"`
}
