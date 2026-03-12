package http

import (
	rbacapp "dispatch/internal/modules/rbac/application"
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, rbacSvc *rbacapp.Service) {
	rg.GET("/logs", rbacmiddleware.RequirePermission(rbacSvc, "fuel.read"), h.List)
	rg.GET("/logs/:id", rbacmiddleware.RequirePermission(rbacSvc, "fuel.read"), h.Get)
	rg.POST("/logs", rbacmiddleware.RequirePermission(rbacSvc, "fuel.manage"), h.Create)
	rg.PUT("/logs/:id", rbacmiddleware.RequirePermission(rbacSvc, "fuel.manage"), h.Update)
	rg.DELETE("/logs/:id", rbacmiddleware.RequirePermission(rbacSvc, "fuel.manage"), h.Delete)
}

