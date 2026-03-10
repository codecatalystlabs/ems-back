package http

import (
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("", rbacmiddleware.RequirePermission(nil, "fleet.read"), h.List)
	rg.GET("/:id", rbacmiddleware.RequirePermission(nil, "fleet.read"), h.Get)
	rg.POST("", rbacmiddleware.RequirePermission(nil, "fleet.manage"), h.Create)
	rg.PUT("/:id", rbacmiddleware.RequirePermission(nil, "fleet.manage"), h.Update)
	rg.DELETE("/:id", rbacmiddleware.RequirePermission(nil, "fleet.manage"), h.Delete)
}

