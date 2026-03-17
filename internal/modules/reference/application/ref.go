package application

import (
	"context"

	"dispatch/internal/modules/reference/application/dto"
	refdomain "dispatch/internal/modules/reference/domain"
	platformdb "dispatch/internal/platform/db"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListDistricts(ctx context.Context, params dto.ListDistrictsParams) (platformdb.PageResult[refdomain.District], error) {
	items, total, err := s.repo.ListDistricts(ctx, params)
	if err != nil {
		return platformdb.PageResult[refdomain.District]{}, err
	}
	return platformdb.PageResult[refdomain.District]{
		Items: items,
		Meta:  platformdb.NewPageMeta(params.Pagination, total),
	}, nil
}

func (s *Service) ListSubcounties(ctx context.Context, params dto.ListSubcountiesParams) (platformdb.PageResult[refdomain.Subcounty], error) {
	items, total, err := s.repo.ListSubcounties(ctx, params)
	if err != nil {
		return platformdb.PageResult[refdomain.Subcounty]{}, err
	}
	return platformdb.PageResult[refdomain.Subcounty]{
		Items: items,
		Meta:  platformdb.NewPageMeta(params.Pagination, total),
	}, nil
}

func (s *Service) ListFacilities(ctx context.Context, params dto.ListFacilitiesParams) (platformdb.PageResult[refdomain.Facility], error) {
	items, total, err := s.repo.ListFacilities(ctx, params)
	if err != nil {
		return platformdb.PageResult[refdomain.Facility]{}, err
	}
	return platformdb.PageResult[refdomain.Facility]{
		Items: items,
		Meta:  platformdb.NewPageMeta(params.Pagination, total),
	}, nil
}

func (s *Service) ListFacilityLevels(ctx context.Context) ([]refdomain.FacilityLevel, error) {
	return s.repo.ListFacilityLevels(ctx)
}

func (s *Service) ListIncidentTypes(ctx context.Context) ([]refdomain.IncidentType, error) {
	return s.repo.ListIncidentTypes(ctx)
}

func (s *Service) ListPriorityLevels(ctx context.Context) ([]refdomain.PriorityLevel, error) {
	return s.repo.ListPriorityLevels(ctx)
}

func (s *Service) ListSeverityLevels(ctx context.Context) ([]refdomain.SeverityLevel, error) {
	return s.repo.ListSeverityLevels(ctx)
}

func (s *Service) ListAmbulanceCategories(ctx context.Context) ([]refdomain.AmbulanceCategory, error) {
	return s.repo.ListAmbulanceCategories(ctx)
}

func (s *Service) ListCapabilities(ctx context.Context) ([]refdomain.Capability, error) {
	return s.repo.ListCapabilities(ctx)
}

func (s *Service) ListTriageQuestions(ctx context.Context, params dto.ListTriageQuestionsParams) (platformdb.PageResult[refdomain.TriageQuestion], error) {
	items, total, err := s.repo.ListTriageQuestions(ctx, params)
	if err != nil {
		return platformdb.PageResult[refdomain.TriageQuestion]{}, err
	}
	return platformdb.PageResult[refdomain.TriageQuestion]{
		Items: items,
		Meta:  platformdb.NewPageMeta(params.Pagination, total),
	}, nil
}

func (s *Service) ListRoles(ctx context.Context, params dto.ListRolesParams) (platformdb.PageResult[refdomain.Role], error) {
	items, total, err := s.repo.ListRoles(ctx, params)
	if err != nil {
		return platformdb.PageResult[refdomain.Role]{}, err
	}
	return platformdb.PageResult[refdomain.Role]{
		Items: items,
		Meta:  platformdb.NewPageMeta(params.Pagination, total),
	}, nil
}
