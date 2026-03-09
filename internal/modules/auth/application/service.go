package application

import "context"

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Login(_ context.Context, username, password string) (string, error) {
	return "implement login", nil
}
