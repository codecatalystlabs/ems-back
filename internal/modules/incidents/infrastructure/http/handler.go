package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	incapp "dispatch/internal/modules/incidents/application"
	platformdb "dispatch/internal/platform/db"
	"dispatch/internal/platform/httpx"
)

type Handler struct {
	service *incapp.Service
}

func NewHandler(service *incapp.Service) *Handler {
	return &Handler{service: service}
}

// List godoc
//
//	@Summary		List incidents
//	@Description	Returns paginated incidents with filters for status and type
//	@Tags			Incidents
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page						query		int		false	"Page number"	default(1)
//	@Param			page_size					query		int		false	"Page size"		default(20)
//	@Param			search						query		string	false	"Search term (number, caller, patient, summary)"
//	@Param			sort_by						query		string	false	"Sort field"	Enums(reported_at,status)
//	@Param			sort_order					query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Param			filter[status]				query		string	false	"Filter by status"				Enums(NEW,PENDING_VERIFICATION,VERIFIED,AWAITING_ASSIGNMENT,ASSIGNED,ENROUTE,AT_SCENE,TRANSPORTING,COMPLETED,CANCELLED,ESCALATED,REJECTED)
//	@Param			filter[verification_status]	query		string	false	"Filter by verification status"	Enums(PENDING,VERIFIED,REJECTED)
//	@Param			filter[incident_type_id]	query		string	false	"Filter by incident type id"
//	@Param			filter[district_id]			query		string	false	"Filter by district id"
//	@Param			filter[facility_id]			query		string	false	"Filter by facility id"
//	@Param			filter[date_from]			query		string	false	"Filter by reported_at from (ISO 8601)"
//	@Param			filter[date_to]				query		string	false	"Filter by reported_at to (ISO 8601)"
//	@Success		200							{object}	map[string]interface{}
//	@Failure		500							{object}	map[string]interface{}
//	@Router			/incidents [get]
func (h *Handler) List(c *gin.Context) {
	p := platformdb.ParsePagination(
		c.Request.URL.Query(),
		map[string]string{
			"created_at":  "i.created_at",
			"reported_at": "i.reported_at",
			"status":      "i.status",
		},
		map[string]struct{}{
			"status":              {},
			"verification_status": {},
			"incident_type_id":    {},
			"district_id":         {},
			"facility_id":         {},
			"date_from":           {},
			"date_to":             {},
		},
	)
	out, err := h.service.List(c.Request.Context(), p)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// Create godoc
//
//	@Summary		Create incident
//	@Description	Creates a new incident record from an emergency alert
//	@Tags			Incidents
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			payload	body		incapp.CreateIncidentRequest	true	"Create incident payload"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/incidents [post]
func (h *Handler) Create(c *gin.Context) {
	var req incapp.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	userID := c.GetString("user_id")
	var createdBy *string
	if userID != "" {
		createdBy = &userID
	}
	out, err := h.service.Create(c.Request.Context(), req, createdBy)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Created(c, out)
}

// Get godoc
//
//	@Summary		Get incident
//	@Description	Get a single incident by ID
//	@Tags			Incidents
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Incident ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/incidents/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	out, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		httpx.Error(c, http.StatusNotFound, err.Error())
		return
	}
	httpx.OK(c, out)
}

// Update godoc
//
//	@Summary		Update incident
//	@Description	Update an existing incident
//	@Tags			Incidents
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Incident ID"
//	@Param			payload	body		incapp.UpdateIncidentRequest	true	"Update incident payload"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/incidents/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req incapp.UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// Delete godoc
//
//	@Summary		Delete incident
//	@Description	Delete an incident by ID
//	@Tags			Incidents
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Incident ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/incidents/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, gin.H{"message": "incident deleted"})
}
