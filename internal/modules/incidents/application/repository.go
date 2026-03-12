package application

import (
	"context"

	"dispatch/internal/modules/incidents/domain"
	platformdb "dispatch/internal/platform/db"
)

type Repository interface {
	ListIncidents(ctx context.Context, p platformdb.Pagination) ([]domain.Incident, int64, error)
	CreateIncident(ctx context.Context, in domain.Incident) (domain.Incident, error)
	GetByID(ctx context.Context, id string) (domain.Incident, error)
	UpdateIncident(ctx context.Context, id string, req UpdateIncidentRequest) (domain.Incident, error)
	DeleteIncident(ctx context.Context, id string) error
}
