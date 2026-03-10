package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	userSv "dispatch/internal/modules/users/application"
	"dispatch/internal/modules/users/application/dto"
	userdto "dispatch/internal/modules/users/application/dto"
	"dispatch/internal/platform/db"
	"dispatch/internal/platform/httpx"
)

type Handler struct {
	service *userSv.Service
}

func NewHandler(service *userSv.Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
//
//	@Summary		Create user
//	@Description	Creates a new system user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			payload	body		userdto.CreateUserRequest	true	"Create user payload"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/users [post]
func (h *Handler) Create(c *gin.Context) {
	var req userdto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Created(c, user)
}

// List godoc
//
//	@Summary		List users
//	@Description	Returns paginated users with search, sorting, and filters
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page				query		int		false	"Page number"	default(1)
//	@Param			page_size			query		int		false	"Page size"		default(20)
//	@Param			search				query		string	false	"Search term"
//	@Param			sort_by				query		string	false	"Sort field"	Enums(created_at,username,first_name,last_name,status)
//	@Param			sort_order			query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Param			filter[status]		query		string	false	"Filter by status"
//	@Param			filter[is_active]	query		string	false	"Filter by active flag"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/users [get]
func (h *Handler) List(c *gin.Context) {
	params := dto.ListUsersParams{
		Pagination: db.ParsePagination(
			c.Request.URL.Query(),
			map[string]string{
				"created_at": "u.created_at",
				"username":   "u.username",
				"first_name": "u.first_name",
				"last_name":  "u.last_name",
				"status":     "u.status",
			},
			map[string]struct{}{
				"status":    {},
				"is_active": {},
			},
		),
	}

	result, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, result)
}

// GetByID godoc
//
//	@Summary		Get user by ID
//	@Description	Returns a user by their ID
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/users/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	user, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, userSv.ErrUserNotFound) {
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		}
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, user)
}

// Update godoc
//
//	@Summary		Update user
//	@Description	Updates user details. All fields are optional.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"User ID"
//	@Param			payload	body		userdto.UpdateUserRequest	true	"Update user payload"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/users/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, userSv.ErrUserNotFound) {
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		}
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, user)
}

// Delete godoc
//
//	@Summary		Delete user
//	@Description	Soft deletes a user by their ID
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, userSv.ErrUserNotFound) {
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		}
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, gin.H{"message": "user deleted"})
}
