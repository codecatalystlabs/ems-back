package application

import (
	"context"

	"dispatch/internal/modules/notifications/domain"
	platformdb "dispatch/internal/platform/db"
)

type Repository interface {
	ListNotifications(ctx context.Context, userID string, p platformdb.Pagination) ([]domain.Notification, int64, error)
	GetByID(ctx context.Context, id string) (domain.Notification, error)
	Create(ctx context.Context, in domain.Notification) (domain.Notification, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
