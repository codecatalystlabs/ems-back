package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMiddleware gin.HandlerFunc) {
	rg.GET("/districts", h.ListDistricts)
	rg.GET("/subcounties", h.ListSubcounties)
	rg.GET("/facilities", h.ListFacilities)

	secured := rg.Group("")
	secured.Use(authMiddleware)
	rg.GET("/facility-levels", h.ListFacilityLevels)
	rg.GET("/incident-types", h.ListIncidentTypes)
	rg.GET("/priority-levels", h.ListPriorityLevels)
	rg.GET("/severity-levels", h.ListSeverityLevels)
	secured.GET("/ambulance-categories", h.ListAmbulanceCategories)
	secured.GET("/capabilities", h.ListCapabilities)

	secured.GET("/triage-questions", h.ListTriageQuestions)
	rg.GET("/roles", h.ListRoles)
}
