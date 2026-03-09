package types

import (
	"dispatch/internal/platform/config"
	"dispatch/internal/platform/events"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ModuleDeps struct {
	Router *gin.RouterGroup
	DB     *pgxpool.Pool
	Redis  *redis.Client
	Logger *zap.Logger
	Bus    events.Publisher
	Config config.Config
}
