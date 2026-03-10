package http

import (
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("", rbacmiddleware.RequirePermission(nil, "incidents.read"), h.List)
	rg.POST("", rbacmiddleware.RequirePermission(nil, "incidents.create"), h.Create)
	rg.GET("/:id", rbacmiddleware.RequirePermission(nil, "incidents.read"), h.Get)
	rg.PUT("/:id", rbacmiddleware.RequirePermission(nil, "incidents.triage"), h.Update)
	rg.DELETE("/:id", rbacmiddleware.RequirePermission(nil, "incidents.triage"), h.Delete)
}

