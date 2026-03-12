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

	rg.POST("/:id/change-password", h.ChangePassword)
	rg.PATCH("/:id/profile", h.UpdateProfile)

	rg.POST("/:id/roles", h.AssignRole)
	rg.DELETE("/:id/roles/:roleId", h.RemoveRole)

	rg.POST("/:id/assignments", h.AssignUser)
	rg.PATCH("/assignments/:assignmentId", h.UpdateAssignment)

	rg.POST("/:id/capabilities", h.AssignCapability)
	rg.PATCH("/capabilities/:capabilityRecordId", h.UpdateCapability)
	rg.GET("/:id/details", h.GetDetails)
}
