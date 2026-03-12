package application

import (
	"context"

	"dispatch/internal/modules/reference/application/dto"

	refdomain "dispatch/internal/modules/reference/domain"
)

type Repository interface {
	ListDistricts(ctx context.Context, params dto.ListDistrictsParams) ([]refdomain.District, int64, error)
	ListSubcounties(ctx context.Context, params dto.ListSubcountiesParams) ([]refdomain.Subcounty, int64, error)
	ListFacilities(ctx context.Context, params dto.ListFacilitiesParams) ([]refdomain.Facility, int64, error)

	ListFacilityLevels(ctx context.Context) ([]refdomain.FacilityLevel, error)
	ListIncidentTypes(ctx context.Context) ([]refdomain.IncidentType, error)
	ListPriorityLevels(ctx context.Context) ([]refdomain.PriorityLevel, error)
	ListSeverityLevels(ctx context.Context) ([]refdomain.SeverityLevel, error)
	ListAmbulanceCategories(ctx context.Context) ([]refdomain.AmbulanceCategory, error)
	ListCapabilities(ctx context.Context) ([]refdomain.Capability, error)
}
