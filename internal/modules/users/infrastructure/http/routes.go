package http

import (
	"github.com/gin-gonic/gin"

	rbacapp "dispatch/internal/modules/rbac/application"
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, rbacSvc *rbacapp.Service) {
	rg.POST("", rbacmiddleware.RequirePermission(rbacSvc, "users.create"), h.Create)
	rg.GET("", rbacmiddleware.RequirePermission(rbacSvc, "users.read"), h.List)
	rg.GET("/:id", rbacmiddleware.RequirePermission(rbacSvc, "users.read"), h.GetByID)
	rg.PUT("/:id", rbacmiddleware.RequirePermission(rbacSvc, "users.update"), h.Update)
	rg.DELETE("/:id", rbacmiddleware.RequirePermission(rbacSvc, "users.delete"), h.Delete)
}
