package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	rbacapp "dispatch/internal/modules/rbac/application"
	"dispatch/internal/platform/httpx"
)

type Handler struct {
	service *rbacapp.Service
}

func NewHandler(service *rbacapp.Service) *Handler { return &Handler{service: service} }

func (h *Handler) MyPermissions(c *gin.Context) {
	userID := c.GetString("user_id")
	items, err := h.service.ListPermissionGrants(c.Request.Context(), userID)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, items)
}
