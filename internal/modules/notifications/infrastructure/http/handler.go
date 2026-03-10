package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	notifapp "dispatch/internal/modules/notifications/application"
	platformdb "dispatch/internal/platform/db"
	"dispatch/internal/platform/httpx"
)

type Handler struct {
	service *notifapp.Service
}

func NewHandler(service *notifapp.Service) *Handler {
	return &Handler{service: service}
}

// ListMy godoc
//
//	@Summary		List my notifications
//	@Tags			Notifications
//	@Produce		json
//	@Param			page		query	int		false	"Page number"	default(1)
//	@Param			page_size	query	int		false	"Page size"		default(20)
//	@Param			filter[status]	query	string	false	"Filter by status"
//	@Param			filter[channel]	query	string	false	"Filter by channel"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/notifications [get]
func (h *Handler) ListMy(c *gin.Context) {
	userID := c.GetString("user_id")
	p := platformdb.ParsePagination(
		c.Request.URL.Query(),
		map[string]string{
			"created_at": "n.created_at",
		},
		map[string]struct{}{
			"status":  {},
			"channel": {},
		},
	)
	out, err := h.service.ListMy(c.Request.Context(), userID, p)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// Get godoc
//
//	@Summary		Get notification
//	@Tags			Notifications
//	@Produce		json
//	@Param			id	path	string	true	"Notification ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/notifications/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	out, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		httpx.Error(c, http.StatusNotFound, err.Error())
		return
	}
	httpx.OK(c, out)
}

// Create godoc
//
//	@Summary		Create notification
//	@Tags			Notifications
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		notifapp.CreateNotificationRequest	true	"Create notification payload"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/notifications [post]
func (h *Handler) Create(c *gin.Context) {
	var req notifapp.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Created(c, out)
}

// MarkRead godoc
//
//	@Summary		Mark notification as read
//	@Tags			Notifications
//	@Produce		json
//	@Param			id	path	string	true	"Notification ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/notifications/{id}/read [post]
func (h *Handler) MarkRead(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.UpdateStatus(c.Request.Context(), id, "READ"); err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, gin.H{"message": "notification marked read"})
}

