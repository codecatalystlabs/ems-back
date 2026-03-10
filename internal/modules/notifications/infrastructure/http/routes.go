package http

import (
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("", rbacmiddleware.RequirePermission(nil, "dispatch.read"), h.ListMy)
	rg.GET("/:id", rbacmiddleware.RequirePermission(nil, "dispatch.read"), h.Get)
	rg.POST("", rbacmiddleware.RequirePermission(nil, "dispatch.read"), h.Create)
	rg.POST("/:id/read", rbacmiddleware.RequirePermission(nil, "dispatch.read"), h.MarkRead)
}

