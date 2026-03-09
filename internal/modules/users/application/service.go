package application

import (
	"context"
	"time"

	"dispatch/internal/platform/auth"
	"dispatch/internal/platform/db"
	"dispatch/internal/platform/events"

	"dispatch/internal/modules/users/application/dto"
	"dispatch/internal/modules/users/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ServiceRepository interface {
	Create(ctx context.Context, user domain.User, passwordHash string) error
	List(ctx context.Context, params dto.ListUsersParams) ([]domain.User, int64, error)
}

type Service struct {
	repo  ServiceRepository
	bus   events.Publisher
	log   *zap.Logger
	topic string
}

func NewService(repo ServiceRepository, bus events.Publisher, log *zap.Logger, topic string) *Service {
	return &Service{repo: repo, bus: bus, log: log, topic: topic}
}

func (s *Service) Create(ctx context.Context, req dto.CreateUserRequest) (domain.User, error) {
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return domain.User{}, err
	}

	u := domain.User{
		ID:        uuid.NewString(),
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    "ACTIVE",
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, u, hash); err != nil {
		return domain.User{}, err
	}

	_ = s.bus.Publish(ctx, s.topic, events.Event{
		ID:          uuid.NewString(),
		Topic:       s.topic,
		AggregateID: u.ID,
		Type:        "user.created",
		OccurredAt:  time.Now().UTC(),
		Payload: map[string]any{
			"user_id":  u.ID,
			"username": u.Username,
		},
	})

	return u, nil
}

func (s *Service) List(ctx context.Context, params dto.ListUsersParams) (db.PageResult[domain.User], error) {
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return db.PageResult[domain.User]{}, err
	}
	return db.PageResult[domain.User]{
		Items: items,
		Meta:  db.NewPageMeta(params.Pagination, total),
	}, nil
}
