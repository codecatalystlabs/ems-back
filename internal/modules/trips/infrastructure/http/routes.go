package http

import (
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("", rbacmiddleware.RequirePermission(nil, "trips.read"), h.List)
	rg.GET("/:id", rbacmiddleware.RequirePermission(nil, "trips.read"), h.Get)
	rg.POST("", rbacmiddleware.RequirePermission(nil, "trips.read"), h.Create)
	rg.PUT("/:id", rbacmiddleware.RequirePermission(nil, "trips.read"), h.Update)
	rg.DELETE("/:id", rbacmiddleware.RequirePermission(nil, "trips.read"), h.Delete)
	rg.GET("/:id/events", rbacmiddleware.RequirePermission(nil, "trips.read"), h.ListEvents)
	rg.POST("/:id/events", rbacmiddleware.RequirePermission(nil, "trips.read"), h.CreateEvent)
}

