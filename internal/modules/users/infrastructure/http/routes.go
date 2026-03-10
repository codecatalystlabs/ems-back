package http

import (
	"github.com/gin-gonic/gin"

	rbacmiddleware "dispatch/internal/modules/rbac/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("", rbacmiddleware.RequirePermission(nil, "users.create"), h.Create)
	rg.GET("", rbacmiddleware.RequirePermission(nil, "users.read"), h.List)
	rg.GET("/:id", rbacmiddleware.RequirePermission(nil, "users.read"), h.GetByID)
	rg.PUT("/:id", rbacmiddleware.RequirePermission(nil, "users.update"), h.Update)
	rg.DELETE("/:id", rbacmiddleware.RequirePermission(nil, "users.delete"), h.Delete)
}
