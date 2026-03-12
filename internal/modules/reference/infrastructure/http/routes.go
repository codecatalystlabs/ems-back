package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/districts", h.ListDistricts)
	rg.GET("/subcounties", h.ListSubcounties)
	rg.GET("/facilities", h.ListFacilities)

	rg.GET("/facility-levels", h.ListFacilityLevels)
	rg.GET("/incident-types", h.ListIncidentTypes)
	rg.GET("/priority-levels", h.ListPriorityLevels)
	rg.GET("/severity-levels", h.ListSeverityLevels)
	rg.GET("/ambulance-categories", h.ListAmbulanceCategories)
	rg.GET("/capabilities", h.ListCapabilities)
}
