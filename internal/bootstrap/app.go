package bootstrap

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"dispatch/internal/platform/config"
	"dispatch/internal/platform/events"
	"dispatch/internal/shared/types"
)

type App struct {
	cfg   config.Config
	log   *zap.Logger
	http  *http.Server
	kafka events.KafkaBus
}

type Worker struct {
	cfg   config.Config
	log   *zap.Logger
	kafka events.KafkaBus
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log, err := NewLogger(cfg.Log)
	if err != nil {
		return nil, err
	}

	db, err := NewPostgres(ctx, cfg.DB)
	if err != nil {
		return nil, err
	}

	redisClient, err := NewRedis(ctx, cfg.Redis)
	if err != nil {
		return nil, err
	}

	r := NewRouter()
	api := r.Group("/api/v1")
	RegisterModules(types.ModuleDeps{
		Router: api,
		DB:     db,
		Redis:  redisClient,
		Logger: log,
		Config: cfg,
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      r,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
	}

	return &App{cfg: cfg, log: log, http: srv}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		a.log.Info("server starting", zap.String("addr", a.http.Addr))
		errCh <- a.http.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.App.ShutdownTimeout)
		defer cancel()
		_ = a.kafka.Close()
		return a.http.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
