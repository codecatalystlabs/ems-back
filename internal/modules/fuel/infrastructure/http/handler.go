package http

import (
	"net/http"

	fuelapp "dispatch/internal/modules/fuel/application"
	"dispatch/internal/modules/fuel/domain"
	platformdb "dispatch/internal/platform/db"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *fuelapp.Service
}

func NewHandler(svc *fuelapp.Service) *Handler {
	return &Handler{svc: svc}
}

// ListFuelLogs godoc
//
//	@Summary		List fuel logs
//	@Description	List fuel logs with pagination
//	@Tags			Fuel
//	@Security		BearerAuth
//	@Param			page				query		int		false	"Page number"			default(1)
//	@Param			page_size			query		int		false	"Page size"				default(20)
//	@Param			search				query		string	false	"Search query"
//	@Param			sort_by				query		string	false	"Sort by field"			default(created_at)
//	@Param			sort_order			query		string	false	"Sort order (ASC/DESC)"	default(DESC)
//	@Param			filter[ambulance_id]	query		string	false	"Filter by ambulance_id (UUID)"
//	@Success		200					{object}	platformdb.PageResult[domain.FuelLog]
//	@Failure		401					{object}	map[string]any
//	@Failure		403					{object}	map[string]any
//	@Failure		500					{object}	map[string]any
//	@Router			/fuel/logs [get]
func (h *Handler) List(c *gin.Context) {
	p := platformdb.ParsePagination(
		c.Request.URL.Query(),
		map[string]string{
			"created_at":  "fl.created_at",
			"filled_at":   "fl.filled_at",
			"liters":      "fl.liters",
			"cost":        "fl.cost",
			"odometer_km": "fl.odometer_km",
		},
		map[string]struct{}{
			"ambulance_id": {},
		},
	)

	items, total, err := h.svc.List(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to list fuel logs"})
		return
	}

	c.JSON(http.StatusOK, platformdb.PageResult[domain.FuelLog]{
		Items: items,
		Meta:  platformdb.NewPageMeta(p, total),
	})
}

// GetFuelLog godoc
//
//	@Summary		Get fuel log by id
//	@Tags			Fuel
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Fuel log ID (UUID)"
//	@Success		200	{object}	map[string]any
//	@Failure		401	{object}	map[string]any
//	@Failure		403	{object}	map[string]any
//	@Failure		404	{object}	map[string]any
//	@Failure		500	{object}	map[string]any
//	@Router			/fuel/logs/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	item, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "fuel log not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateFuelLog godoc
//
//	@Summary		Create fuel log
//	@Tags			Fuel
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		fuelapp.CreateFuelLogRequest	true	"Fuel log payload"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	map[string]any
//	@Failure		401		{object}	map[string]any
//	@Failure		403		{object}	map[string]any
//	@Failure		500		{object}	map[string]any
//	@Router			/fuel/logs [post]
func (h *Handler) Create(c *gin.Context) {
	var req fuelapp.CreateFuelLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	var filledBy *string
	if v := c.GetString("user_id"); v != "" {
		filledBy = &v
	}
	item, err := h.svc.Create(c.Request.Context(), req, filledBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create fuel log"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateFuelLog godoc
//
//	@Summary		Update fuel log
//	@Tags			Fuel
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Fuel log ID (UUID)"
//	@Param			payload	body		fuelapp.UpdateFuelLogRequest	true	"Update payload"
//	@Success		200		{object}	map[string]any
//	@Failure		400		{object}	map[string]any
//	@Failure		401		{object}	map[string]any
//	@Failure		403		{object}	map[string]any
//	@Failure		404		{object}	map[string]any
//	@Failure		500		{object}	map[string]any
//	@Router			/fuel/logs/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req fuelapp.UpdateFuelLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	item, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "fuel log not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteFuelLog godoc
//
//	@Summary		Delete fuel log
//	@Tags			Fuel
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Fuel log ID (UUID)"
//	@Success		204	{object}	nil
//	@Failure		401	{object}	map[string]any
//	@Failure		403	{object}	map[string]any
//	@Failure		404	{object}	map[string]any
//	@Failure		500	{object}	map[string]any
//	@Router			/fuel/logs/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "fuel log not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

