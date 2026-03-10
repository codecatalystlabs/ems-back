package application

import (
	"context"

	"dispatch/internal/modules/fleet/domain"
	platformdb "dispatch/internal/platform/db"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(repo Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) ListAmbulances(ctx context.Context, p platformdb.Pagination) (platformdb.PageResult[domain.Ambulance], error) {
	items, total, err := s.repo.ListAmbulances(ctx, p)
	if err != nil {
		return platformdb.PageResult[domain.Ambulance]{}, err
	}
	return platformdb.PageResult[domain.Ambulance]{
		Items: items,
		Meta:  platformdb.NewPageMeta(p, total),
	}, nil
}

func (s *Service) GetAmbulance(ctx context.Context, id string) (domain.Ambulance, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) CreateAmbulance(ctx context.Context, req CreateAmbulanceRequest) (domain.Ambulance, error) {
	a := domain.Ambulance{
		Code:              req.Code,
		PlateNumber:       req.PlateNumber,
		VIN:               req.VIN,
		Make:              req.Make,
		Model:             req.Model,
		YearOfManufacture: req.YearOfManufacture,
		CategoryID:        req.CategoryID,
		OwnershipType:     req.OwnershipType,
		StationFacilityID: req.StationFacilityID,
		DistrictID:        req.DistrictID,
		Status:            "AVAILABLE",
		DispatchReadiness: "DISPATCHABLE",
		IsActive:          true,
	}
	if req.Status != nil {
		a.Status = *req.Status
	}
	if req.DispatchReadiness != nil {
		a.DispatchReadiness = *req.DispatchReadiness
	}
	return s.repo.Create(ctx, a)
}

func (s *Service) UpdateAmbulance(ctx context.Context, id string, req UpdateAmbulanceRequest) (domain.Ambulance, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *Service) DeleteAmbulance(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}


