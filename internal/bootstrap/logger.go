package bootstrap

import (
	"dispatch/internal/platform/config"

	"go.uber.org/zap"
)

func NewLogger(cfg config.LogConfig) (*zap.Logger, error) {
	if cfg.Level == "debug" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
