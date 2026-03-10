package http

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authapp "dispatch/internal/modules/auth/application"
	dto "dispatch/internal/modules/auth/application/dto"
	"dispatch/internal/platform/httpx"
)

type Handler struct {
	service *authapp.Service
}

func NewHandler(service *authapp.Service) *Handler {
	return &Handler{service: service}
}

// Login godoc
//
//	@Summary		Login
//	@Description	Authenticates a user and returns access and refresh tokens
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.LoginRequest	true	"Login payload"
//	@Success		200		{object}	dto.AuthResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	deviceID := c.GetString("device_id")
	deviceName := c.GetString("device_name")
	out, err := h.service.Login(c.Request.Context(), req, deviceID, deviceName, clientIP(c.ClientIP()), c.Request.UserAgent())
	if err != nil {
		switch {
		case errors.Is(err, authapp.ErrInvalidCredentials):
			httpx.Error(c, http.StatusUnauthorized, "invalid credentials")
		case errors.Is(err, authapp.ErrInactiveUser):
			httpx.Error(c, http.StatusForbidden, "user inactive")
		case errors.Is(err, authapp.ErrLockedUser):
			httpx.Error(c, http.StatusForbidden, "user locked")
		default:
			httpx.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	httpx.OK(c, out)
}

// Refresh godoc
//
//	@Summary		Refresh tokens
//	@Description	Refreshes the access token using the refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.RefreshRequest	true	"Refresh payload"
//	@Success		200		{object}	dto.AuthResponse
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		httpx.Error(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	httpx.OK(c, out)
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Revokes current session or all sessions
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			payload	body		dto.LogoutRequest	true	"Logout payload"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	userID := c.GetString("user_id")
	var err error
	if req.LogoutAll {
		err = h.service.LogoutAll(c.Request.Context(), userID)
	} else {
		err = h.service.Logout(c.Request.Context(), req.RefreshToken)
	}
	if err != nil {
		httpx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	httpx.OK(c, gin.H{"message": "logged out"})
}

// Sessions godoc
//
//	@Summary		List sessions
//	@Description	Lists active user sessions
//	@Tags			Auth
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}
//	@Router			/auth/sessions [get]
func (h *Handler) Sessions(c *gin.Context) {
	userID := c.GetString("user_id")
	items, err := h.service.Sessions(c.Request.Context(), userID)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, items)
}

func clientIP(raw string) string {
	ip := strings.TrimSpace(raw)
	if host, _, err := net.SplitHostPort(ip); err == nil {
		return host
	}
	return ip
}
