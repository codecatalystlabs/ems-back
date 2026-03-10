package application

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	rbacdomain "dispatch/internal/modules/rbac/domain"
)

type Service struct {
	repo  Repository
	redis *redis.Client
	log   *zap.Logger
}

var ErrRBACNotInitialized = errors.New("rbac service not initialized")

func NewService(repo Repository, redis *redis.Client, log *zap.Logger) *Service {
	return &Service{repo: repo, redis: redis, log: log}
}

func (s *Service) ListPermissionGrants(ctx context.Context, userID string) ([]rbacdomain.PermissionGrant, error) {
	if s == nil || s.repo == nil {
		return nil, ErrRBACNotInitialized
	}
	return s.repo.ListPermissionGrants(ctx, userID)
}

func (s *Service) HasPermission(ctx context.Context, userID, permission, scopeType string, scopeID *string) (bool, error) {
	if s == nil || s.repo == nil {
		return false, ErrRBACNotInitialized
	}

	grants, err := s.repo.ListPermissionGrants(ctx, userID)
	if err != nil {
		return false, err
	}

	permission = normalize(permission)
	scopeType = normalize(scopeType)

	for _, g := range grants {
		if !permissionMatches(normalize(g.PermCode), permission) {
			continue
		}
		if scopeAllowed(g.ScopeType, g.ScopeID, scopeType, scopeID) {
			return true, nil
		}
	}

	return false, nil
}

func normalize(v string) string {
	return strings.ToUpper(strings.TrimSpace(v))
}

func permissionMatches(grant, required string) bool {
	if grant == required {
		return true
	}
	if grant == "*" || grant == "*.*" {
		return true
	}
	if strings.HasSuffix(grant, ".*") {
		prefix := strings.TrimSuffix(grant, ".*")
		return strings.HasPrefix(required, prefix+".")
	}
	return false
}

func scopeAllowed(grantScopeType string, grantScopeID *string, reqScopeType string, reqScopeID *string) bool {
	gst := normalize(grantScopeType)
	rst := normalize(reqScopeType)

	if gst == "GLOBAL" || rst == "" {
		return true
	}
	if gst != rst {
		return false
	}
	if grantScopeID == nil || reqScopeID == nil {
		return true
	}
	return *grantScopeID == *reqScopeID
}
