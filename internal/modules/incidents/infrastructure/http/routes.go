package http

import (
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"

	"github.com/gin-gonic/gin"

	rbacapp "dispatch/internal/modules/rbac/application"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, rbacSvc *rbacapp.Service) {
	rg.GET("", rbacmiddleware.RequirePermission(rbacSvc, "incidents.read"), h.List)
	rg.POST("", rbacmiddleware.RequirePermission(rbacSvc, "incidents.create"), h.Create)
	rg.GET("/:id", rbacmiddleware.RequirePermission(rbacSvc, "incidents.read"), h.GetByID)
	rg.PATCH("/:id/status", rbacmiddleware.RequirePermission(rbacSvc, "incidents.triage"), h.UpdateStatus)
}
