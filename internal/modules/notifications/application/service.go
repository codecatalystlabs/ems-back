package application

import (
	"context"
	"strings"
	"time"

	"dispatch/internal/modules/notifications/domain"
	platformdb "dispatch/internal/platform/db"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(repo Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) ListMy(ctx context.Context, userID string, p platformdb.Pagination) (platformdb.PageResult[domain.Notification], error) {
	items, total, err := s.repo.ListNotifications(ctx, userID, p)
	if err != nil {
		return platformdb.PageResult[domain.Notification]{}, err
	}
	return platformdb.PageResult[domain.Notification]{Items: items, Meta: platformdb.NewPageMeta(p, total)}, nil
}

func (s *Service) Get(ctx context.Context, id string) (domain.Notification, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateNotificationRequest) (domain.Notification, error) {
	now := time.Now().UTC()
	n := domain.Notification{
		ID:              uuid.NewString(),
		Type:            req.Type,
		RecipientUserID: req.RecipientUserID,
		RecipientPhone:  req.RecipientPhone,
		RecipientEmail:  req.RecipientEmail,
		Title:           req.Title,
		Body:            req.Body,
		Channel:         strings.ToUpper(strings.TrimSpace(req.Channel)),
		LinkedEntityType: req.LinkedEntityType,
		LinkedEntityID:  req.LinkedEntityID,
		Status:          "PENDING",
		Attempts:        0,
		CreatedAt:       now,
	}
	return s.repo.Create(ctx, n)
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status string) error {
	return s.repo.UpdateStatus(ctx, id, strings.ToUpper(strings.TrimSpace(status)))
}

