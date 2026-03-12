package application

import (
	"context"

	"dispatch/internal/modules/fleet/domain"
	platformdb "dispatch/internal/platform/db"
)

type Repository interface {
	ListAmbulances(ctx context.Context, p platformdb.Pagination) ([]domain.Ambulance, int64, error)
	GetByID(ctx context.Context, id string) (domain.Ambulance, error)
	Create(ctx context.Context, in domain.Ambulance) (domain.Ambulance, error)
	Update(ctx context.Context, id string, req UpdateAmbulanceRequest) (domain.Ambulance, error)
	Delete(ctx context.Context, id string) error
}
