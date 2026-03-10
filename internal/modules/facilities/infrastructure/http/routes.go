package http

import (
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("", rbacmiddleware.RequirePermission(nil, "facilities.read"), h.List)
	rg.GET("/:uid", rbacmiddleware.RequirePermission(nil, "facilities.read"), h.Get)
	rg.POST("", rbacmiddleware.RequirePermission(nil, "facilities.read"), h.Create)
	rg.PUT("/:uid", rbacmiddleware.RequirePermission(nil, "facilities.read"), h.Update)
	rg.DELETE("/:uid", rbacmiddleware.RequirePermission(nil, "facilities.read"), h.Delete)
}

