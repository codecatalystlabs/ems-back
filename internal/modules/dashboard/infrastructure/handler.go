package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dashboardapp "dispatch/internal/modules/dashboard/application"
	"dispatch/internal/platform/httpx"
)

type Handler struct {
	service *dashboardapp.Service
}

func NewHandler(service *dashboardapp.Service) *Handler { return &Handler{service: service} }

// GetEMSOverview godoc
// @Summary EMS dashboard overview
// @Description Returns the EMS dashboard metrics, charts, and ambulance status table based on date and location filters
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Param date_from query string false "Start date YYYY-MM-DD"
// @Param date_to query string false "End date YYYY-MM-DD"
// @Param district_id query string false "District ID"
// @Param facility_id query string false "Facility ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dashboard/ems-overview [get]
func (h *Handler) GetEMSOverview(c *gin.Context) {
	var q dashboardapp.DashboardQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.GetDashboard(c.Request.Context(), q)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}
