package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	availabilityapp "dispatch/internal/modules/availability/application"
	availdto "dispatch/internal/modules/availability/application/dto"
	platformdb "dispatch/internal/platform/db"
	"dispatch/internal/platform/httpx"
)

type Handler struct{ service *availabilityapp.Service }

func NewHandler(service *availabilityapp.Service) *Handler { return &Handler{service: service} }

// CreateShift godoc
//
//	@Summary		Create user shift
//	@Description	Creates a new shift for a user
//	@Tags			Availability
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			payload	body		availdto.CreateShiftRequest	true	"Create shift payload"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/availability/shifts [post]
func (h *Handler) CreateShift(c *gin.Context) {
	var req availdto.CreateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.CreateShift(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Created(c, out)
}

// UpdateShift godoc
//
//	@Summary		Update user shift
//	@Description	Partially updates a shift by ID
//	@Tags			Availability
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"Shift ID"
//	@Param			payload	body		availdto.UpdateShiftRequest	true	"Update shift payload"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/availability/shifts/{id} [patch]
func (h *Handler) UpdateShift(c *gin.Context) {
	var req availdto.UpdateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.UpdateShift(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// GetShiftByID godoc
//
//	@Summary		Get shift by ID
//	@Description	Returns a shift by ID
//	@Tags			Availability
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Shift ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/availability/shifts/{id} [get]
func (h *Handler) GetShiftByID(c *gin.Context) {
	out, err := h.service.GetShiftByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListShifts godoc
//
//	@Summary		List user shifts
//	@Description	Returns paginated shifts with optional filters
//	@Tags			Availability
//	@Produce		json
//	@Security		BearerAuth
//	@Param			user_id		query		string	false	"User ID"
//	@Param			district_id	query		string	false	"District ID"
//	@Param			facility_id	query		string	false	"Facility ID"
//	@Param			shift_date	query		string	false	"Shift date (YYYY-MM-DD)"
//	@Param			status		query		string	false	"Shift status"
//	@Param			page		query		int		false	"Page number"	default(1)
//	@Param			page_size	query		int		false	"Page size"		default(20)
//	@Param			search		query		string	false	"Search shift type"
//	@Param			sort_by		query		string	false	"Sort field"	Enums(shift_date,starts_at,status,created_at)
//	@Param			sort_order	query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Success		200			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]interface{}
//	@Router			/availability/shifts [get]
func (h *Handler) ListShifts(c *gin.Context) {
	var userID, districtID, facilityID, shiftDate, status *string
	if v := c.Query("user_id"); v != "" {
		userID = &v
	}
	if v := c.Query("district_id"); v != "" {
		districtID = &v
	}
	if v := c.Query("facility_id"); v != "" {
		facilityID = &v
	}
	if v := c.Query("shift_date"); v != "" {
		shiftDate = &v
	}
	if v := c.Query("status"); v != "" {
		status = &v
	}
	params := availdto.ListShiftsParams{UserID: userID, DistrictID: districtID, FacilityID: facilityID, ShiftDate: shiftDate, Status: status,
		Pagination: platformdb.ParsePagination(c.Request.URL.Query(), map[string]string{"shift_date": "us.shift_date", "starts_at": "us.starts_at", "status": "us.status", "created_at": "us.created_at"}, map[string]struct{}{}),
	}
	out, err := h.service.ListShifts(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// UpsertAvailability godoc
//
//	@Summary		Upsert user availability
//	@Description	Creates or updates a user's availability record
//	@Tags			Availability
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			payload	body		availdto.UpsertAvailabilityRequest	true	"Availability payload"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/availability/users/availability [post]
func (h *Handler) UpsertAvailability(c *gin.Context) {
	var req availdto.UpsertAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.UpsertAvailability(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// GetAvailabilityByUserID godoc
//
//	@Summary		Get user availability
//	@Description	Returns availability for a specific user
//	@Tags			Availability
//	@Produce		json
//	@Security		BearerAuth
//	@Param			userId	path		string	true	"User ID"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/availability/users/{userId}/availability [get]
func (h *Handler) GetAvailabilityByUserID(c *gin.Context) {
	out, err := h.service.GetAvailabilityByUserID(c.Request.Context(), c.Param("userId"))
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListAvailability godoc
//
//	@Summary		List user availability
//	@Description	Returns paginated user availability records
//	@Tags			Availability
//	@Produce		json
//	@Security		BearerAuth
//	@Param			status			query		string	false	"Availability status"
//	@Param			dispatchable	query		bool	false	"Dispatchable flag"
//	@Param			page			query		int		false	"Page number"	default(1)
//	@Param			page_size		query		int		false	"Page size"		default(20)
//	@Param			sort_by			query		string	false	"Sort field"	Enums(updated_at,availability_status,last_seen_at)
//	@Param			sort_order		query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Success		200				{object}	map[string]interface{}
//	@Failure		500				{object}	map[string]interface{}
//	@Router			/availability/users/availability [get]
func (h *Handler) ListAvailability(c *gin.Context) {
	var status *string
	var dispatchable *bool
	if v := c.Query("status"); v != "" {
		status = &v
	}
	if v := c.Query("dispatchable"); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			dispatchable = &b
		}
	}
	params := availdto.ListAvailabilityParams{
		Status: status, Dispatchable: dispatchable,
		Pagination: platformdb.ParsePagination(c.Request.URL.Query(), map[string]string{"updated_at": "ua.updated_at", "availability_status": "ua.availability_status", "last_seen_at": "ua.last_seen_at"}, map[string]struct{}{}),
	}
	out, err := h.service.ListAvailability(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// CreatePresenceLog godoc
//
//	@Summary		Create presence log
//	@Description	Stores a user presence log entry
//	@Tags			Availability
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			payload	body		availdto.CreatePresenceLogRequest	true	"Presence log payload"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/availability/users/presence-logs [post]
func (h *Handler) CreatePresenceLog(c *gin.Context) {
	var req availdto.CreatePresenceLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.CreatePresenceLog(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Created(c, out)
}

// ListPresenceLogs godoc
//
//	@Summary		List presence logs
//	@Description	Returns paginated user presence logs
//	@Tags			Availability
//	@Produce		json
//	@Security		BearerAuth
//	@Param			user_id		query		string	false	"User ID"
//	@Param			channel		query		string	false	"Channel"
//	@Param			page		query		int		false	"Page number"	default(1)
//	@Param			page_size	query		int		false	"Page size"		default(20)
//	@Param			sort_by		query		string	false	"Sort field"	Enums(seen_at,channel)
//	@Param			sort_order	query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Success		200			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]interface{}
//	@Router			/availability/users/presence-logs [get]
func (h *Handler) ListPresenceLogs(c *gin.Context) {
	var userID, channel *string
	if v := c.Query("user_id"); v != "" {
		userID = &v
	}
	if v := c.Query("channel"); v != "" {
		channel = &v
	}
	params := availdto.ListPresenceParams{UserID: userID, Channel: channel,
		Pagination: platformdb.ParsePagination(c.Request.URL.Query(), map[string]string{"seen_at": "upl.seen_at", "channel": "upl.channel"}, map[string]struct{}{}),
	}
	out, err := h.service.ListPresenceLogs(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}
