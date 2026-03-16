package infrastructure

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	deviceapp "dispatch/internal/modules/device_tokens/application"
	platformdb "dispatch/internal/platform/db"
	"dispatch/internal/platform/httpx"
)

type Handler struct {
	service *deviceapp.Service
}

func NewHandler(service *deviceapp.Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
//
//	@Summary		Register device token
//	@Description	Registers or reactivates a push token for a user
//	@Tags			DeviceTokens
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			payload	body		deviceapp.RegisterDeviceTokenRequest	true	"Device token payload"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/device-tokens [post]
func (h *Handler) Register(c *gin.Context) {
	var req deviceapp.RegisterDeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Created(c, out)
}

// Update godoc
//
//	@Summary		Update device token
//	@Description	Updates a device token record by ID
//	@Tags			DeviceTokens
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string									true	"Device token ID"
//	@Param			payload	body		deviceapp.UpdateDeviceTokenRequest	true	"Update payload"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/device-tokens/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	var req deviceapp.UpdateDeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, deviceapp.ErrDeviceTokenNotFound) {
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		}
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// Deactivate godoc
//
//	@Summary		Deactivate device token
//	@Description	Deactivates a device token by ID
//	@Tags			DeviceTokens
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Device token ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/device-tokens/{id} [delete]
func (h *Handler) Deactivate(c *gin.Context) {
	err := h.service.Deactivate(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, deviceapp.ErrDeviceTokenNotFound) {
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		}
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, gin.H{"message": "device token deactivated"})
}

// GetByID godoc
//
//	@Summary		Get device token by ID
//	@Description	Returns a device token by ID
//	@Tags			DeviceTokens
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Device token ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/device-tokens/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	out, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, deviceapp.ErrDeviceTokenNotFound) {
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		}
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// List godoc
//
//	@Summary		List device tokens
//	@Description	Returns paginated device tokens
//	@Tags			DeviceTokens
//	@Produce		json
//	@Security		BearerAuth
//	@Param			user_id		query		string	false	"User ID"
//	@Param			platform	query		string	false	"Platform"
//	@Param			is_active	query		bool	false	"Active flag"
//	@Param			page		query		int		false	"Page number"	default(1)
//	@Param			page_size	query		int		false	"Page size"		default(20)
//	@Param			search		query		string	false	"Search by device ID or token"
//	@Param			sort_by		query		string	false	"Sort field"	Enums(created_at,updated_at,platform,last_seen_at)
//	@Param			sort_order	query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Success		200			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]interface{}
//	@Router			/device-tokens [get]
func (h *Handler) List(c *gin.Context) {
	var userID, platform *string
	var isActive *bool

	if v := c.Query("user_id"); v != "" {
		userID = &v
	}
	if v := c.Query("platform"); v != "" {
		platform = &v
	}
	if v := c.Query("is_active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			isActive = &b
		}
	}

	params := deviceapp.ListDeviceTokensParams{
		UserID:   userID,
		Platform: platform,
		IsActive: isActive,
		Pagination: platformdb.ParsePagination(
			c.Request.URL.Query(),
			map[string]string{
				"created_at":   "udt.created_at",
				"updated_at":   "udt.updated_at",
				"platform":     "udt.platform",
				"last_seen_at": "udt.last_seen_at",
			},
			map[string]struct{}{},
		),
	}

	out, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}
