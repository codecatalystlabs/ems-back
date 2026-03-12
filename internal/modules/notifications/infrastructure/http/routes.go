package http

import (
	rbacapp "dispatch/internal/modules/rbac/application"
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, rbacSvc *rbacapp.Service) {
	rg.GET("", rbacmiddleware.RequirePermission(rbacSvc, "notifications.read"), h.ListMy)
	rg.GET("/:id", rbacmiddleware.RequirePermission(rbacSvc, "notifications.read"), h.Get)
	rg.POST("", rbacmiddleware.RequirePermission(rbacSvc, "notifications.manage"), h.Create)
	rg.POST("/:id/read", rbacmiddleware.RequirePermission(rbacSvc, "notifications.manage"), h.MarkRead)
}

