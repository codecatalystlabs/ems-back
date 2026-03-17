package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	refapp "dispatch/internal/modules/reference/application"
	"dispatch/internal/modules/reference/application/dto"
	platformdb "dispatch/internal/platform/db"
	"dispatch/internal/platform/httpx"
)

type Handler struct {
	service *refapp.Service
}

func NewHandler(service *refapp.Service) *Handler {
	return &Handler{service: service}
}

// ListDistricts godoc
//
//	@Summary		List districts
//	@Description	Returns paginated districts
//	@Tags			Reference
//	@Produce		json
//	@Param			page				query		int		false	"Page number"	default(1)
//	@Param			page_size			query		int		false	"Page size"		default(20)
//	@Param			search				query		string	false	"Search by district name, code, or region"
//	@Param			sort_by				query		string	false	"Sort field"	Enums(name,region,created_at)
//	@Param			sort_order			query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Param			filter[is_active]	query		string	false	"Filter by active flag"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/reference/districts [get]
func (h *Handler) ListDistricts(c *gin.Context) {
	params := dto.ListDistrictsParams{
		Pagination: platformdb.ParsePagination(
			c.Request.URL.Query(),
			map[string]string{
				"name":       "d.name",
				"region":     "d.region",
				"created_at": "d.created_at",
			},
			map[string]struct{}{
				"is_active": {},
			},
		),
	}

	out, err := h.service.ListDistricts(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListSubcounties godoc
//
//	@Summary		List subcounties
//	@Description	Returns paginated subcounties, optionally filtered by district
//	@Tags			Reference
//	@Produce		json
//	@Param			district_id			query		string	false	"District ID"
//	@Param			page				query		int		false	"Page number"	default(1)
//	@Param			page_size			query		int		false	"Page size"		default(20)
//	@Param			search				query		string	false	"Search by subcounty name or code"
//	@Param			sort_by				query		string	false	"Sort field"	Enums(name,created_at)
//	@Param			sort_order			query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Param			filter[is_active]	query		string	false	"Filter by active flag"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		500					{object}	map[string]interface{}
//	@Router			/reference/subcounties [get]
func (h *Handler) ListSubcounties(c *gin.Context) {
	districtID := c.Query("district_id")
	var districtIDPtr *string
	if districtID != "" {
		districtIDPtr = &districtID
	}

	params := dto.ListSubcountiesParams{
		DistrictID: districtIDPtr,
		Pagination: platformdb.ParsePagination(
			c.Request.URL.Query(),
			map[string]string{
				"name":       "s.name",
				"created_at": "s.created_at",
			},
			map[string]struct{}{
				"is_active": {},
			},
		),
	}

	out, err := h.service.ListSubcounties(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListFacilities godoc
//
//	@Summary		List facilities
//	@Description	Returns paginated facilities, optionally filtered by district, subcounty, or facility level
//	@Tags			Reference
//	@Produce		json
//	@Param			district_id					query		string	false	"District ID"
//	@Param			subcounty_id				query		string	false	"Subcounty ID"
//	@Param			level_id					query		string	false	"Facility level ID"
//	@Param			page						query		int		false	"Page number"	default(1)
//	@Param			page_size					query		int		false	"Page size"		default(20)
//	@Param			search						query		string	false	"Search by facility name, short name, NHFR ID, or code"
//	@Param			sort_by						query		string	false	"Sort field"	Enums(name,created_at,ownership)
//	@Param			sort_order					query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Param			filter[is_active]			query		string	false	"Filter by active flag"
//	@Param			filter[is_dispatch_station]	query		string	false	"Filter by dispatch station flag"
//	@Success		200							{object}	map[string]interface{}
//	@Failure		500							{object}	map[string]interface{}
//	@Router			/reference/facilities [get]
func (h *Handler) ListFacilities(c *gin.Context) {
	districtID := c.Query("district_id")
	subcountyID := c.Query("subcounty_id")
	levelID := c.Query("level_id")

	var districtIDPtr, subcountyIDPtr, levelIDPtr *string
	if districtID != "" {
		districtIDPtr = &districtID
	}
	if subcountyID != "" {
		subcountyIDPtr = &subcountyID
	}
	if levelID != "" {
		levelIDPtr = &levelID
	}

	params := dto.ListFacilitiesParams{
		DistrictID:  districtIDPtr,
		SubcountyID: subcountyIDPtr,
		LevelID:     levelIDPtr,
		Pagination: platformdb.ParsePagination(
			c.Request.URL.Query(),
			map[string]string{
				"name":       "f.name",
				"created_at": "f.created_at",
				"ownership":  "f.ownership",
			},
			map[string]struct{}{
				"is_active":           {},
				"is_dispatch_station": {},
			},
		),
	}

	out, err := h.service.ListFacilities(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListFacilityLevels godoc
//
//	@Summary		List facility levels
//	@Description	Returns active facility levels
//	@Tags			Reference
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/reference/facility-levels [get]
func (h *Handler) ListFacilityLevels(c *gin.Context) {
	out, err := h.service.ListFacilityLevels(c.Request.Context())
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListIncidentTypes godoc
//
//	@Summary		List incident types
//	@Description	Returns active incident types
//	@Tags			Reference
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/reference/incident-types [get]
func (h *Handler) ListIncidentTypes(c *gin.Context) {
	out, err := h.service.ListIncidentTypes(c.Request.Context())
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListPriorityLevels godoc
//
//	@Summary		List priority levels
//	@Description	Returns active priority levels
//	@Tags			Reference
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/reference/priority-levels [get]
func (h *Handler) ListPriorityLevels(c *gin.Context) {
	out, err := h.service.ListPriorityLevels(c.Request.Context())
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListSeverityLevels godoc
//
//	@Summary		List severity levels
//	@Description	Returns active severity levels
//	@Tags			Reference
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/reference/severity-levels [get]
func (h *Handler) ListSeverityLevels(c *gin.Context) {
	out, err := h.service.ListSeverityLevels(c.Request.Context())
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListAmbulanceCategories godoc
//
//	@Summary		List ambulance categories
//	@Description	Returns active ambulance categories
//	@Tags			Reference
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/reference/ambulance-categories [get]
func (h *Handler) ListAmbulanceCategories(c *gin.Context) {
	out, err := h.service.ListAmbulanceCategories(c.Request.Context())
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListCapabilities godoc
//
//	@Summary		List capabilities
//	@Description	Returns active capabilities
//	@Tags			Reference
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/reference/capabilities [get]
func (h *Handler) ListCapabilities(c *gin.Context) {
	out, err := h.service.ListCapabilities(c.Request.Context())
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListTriageQuestions godoc
//
//	@Summary		List triage questions
//	@Description	Returns paginated triage questions, optionally filtered by questionnaire code
//	@Tags			Reference
//	@Produce		json
//	@Security		BearerAuth
//	@Param			questionnaire_code		query		string	false	"Questionnaire code"
//	@Param			page					query		int		false	"Page number"	default(1)
//	@Param			page_size				query		int		false	"Page size"		default(20)
//	@Param			search					query		string	false	"Search by question code, text, or response type"
//	@Param			sort_by					query		string	false	"Sort field"	Enums(display_order,created_at,code)
//	@Param			sort_order				query		string	false	"Sort order"	Enums(ASC,DESC)
//	@Param			filter[is_active]		query		string	false	"Filter by active flag"
//	@Param			filter[is_required]		query		string	false	"Filter by required flag"
//	@Param			filter[response_type]	query		string	false	"Filter by response type"
//	@Success		200						{object}	map[string]interface{}
//	@Failure		500						{object}	map[string]interface{}
//	@Router			/reference/triage-questions [get]
func (h *Handler) ListTriageQuestions(c *gin.Context) {
	questionnaireCode := c.Query("questionnaire_code")
	var questionnaireCodePtr *string
	if questionnaireCode != "" {
		questionnaireCodePtr = &questionnaireCode
	}

	params := dto.ListTriageQuestionsParams{
		QuestionnaireCode: questionnaireCodePtr,
		Pagination: platformdb.ParsePagination(
			c.Request.URL.Query(),
			map[string]string{
				"display_order": "tq.display_order",
				"created_at":    "tq.created_at",
				"code":          "tq.code",
			},
			map[string]struct{}{
				"is_active":     {},
				"is_required":   {},
				"response_type": {},
			},
		),
	}

	out, err := h.service.ListTriageQuestions(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

// ListRoles godoc
// @Summary List roles
// @Description Returns paginated roles (RBAC reference data)
// @Tags Reference
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search by name/code"
// @Param sort_by query string false "Sort field" Enums(name,code,created_at)
// @Param sort_order query string false "Sort order" Enums(ASC,DESC)
// @Param filter[is_system] query string false "Filter system roles"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reference/roles [get]
func (h *Handler) ListRoles(c *gin.Context) {
	params := dto.ListRolesParams{
		Pagination: platformdb.ParsePagination(
			c.Request.URL.Query(),
			map[string]string{
				"name":       "r.name",
				"code":       "r.code",
				"created_at": "r.created_at",
			},
			map[string]struct{}{
				"is_system": {},
			},
		),
	}

	out, err := h.service.ListRoles(c.Request.Context(), params)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.OK(c, out)
}
