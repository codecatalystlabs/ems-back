package application

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshService struct {
	db          *pgxpool.Pool
	mu          sync.Mutex
	lastRefresh time.Time
	minInterval time.Duration
}

func NewRefreshService(db *pgxpool.Pool) *RefreshService {
	return &RefreshService{
		db:          db,
		minInterval: 1 * time.Minute,
	}
}

func (s *RefreshService) RefreshAll(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if time.Since(s.lastRefresh) < s.minInterval {
		return nil
	}

	if _, err := s.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY mv_ambulance_latest_readiness`); err != nil {
		return err
	}
	if _, err := s.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY mv_incident_daily_stats`); err != nil {
		return err
	}
	if _, err := s.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY mv_dashboard_daily_summary`); err != nil {
		return err
	}

	s.lastRefresh = time.Now().UTC()
	return nil
}
