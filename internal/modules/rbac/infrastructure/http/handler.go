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

// MyPermissions godoc
//
//	@Summary		Get my permissions
//	@Description	Returns the list of permissions granted to the authenticated user
//	@Tags			RBAC
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/rbac/me/permissions [get]
func (h *Handler) MyPermissions(c *gin.Context) {
	userID := c.GetString("user_id")
	items, err := h.service.ListPermissionGrants(c.Request.Context(), userID)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, items)
}
