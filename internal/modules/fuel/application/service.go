package application

import (
	"context"
	"time"

	"dispatch/internal/modules/fuel/domain"
	platformdb "dispatch/internal/platform/db"

	"go.uber.org/zap"
)

type Service struct {
	repo   Repository
	log    *zap.Logger
}

func NewService(repo Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) List(ctx context.Context, p platformdb.Pagination) ([]domain.FuelLog, int64, error) {
	return s.repo.List(ctx, p)
}

func (s *Service) Get(ctx context.Context, id string) (domain.FuelLog, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateFuelLogRequest, filledByUserID *string) (domain.FuelLog, error) {
	now := time.Now()
	filledAt := now
	if req.FilledAt != nil {
		filledAt = *req.FilledAt
	}
	in := domain.FuelLog{
		AmbulanceID:  req.AmbulanceID,
		FuelType:     req.FuelType,
		Liters:       req.Liters,
		Cost:         req.Cost,
		OdometerKM:   req.OdometerKM,
		StationName:  req.StationName,
		FilledAt:     filledAt,
		FilledBy:     filledByUserID,
		Notes:        req.Notes,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return s.repo.Create(ctx, in)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateFuelLogRequest) (domain.FuelLog, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

