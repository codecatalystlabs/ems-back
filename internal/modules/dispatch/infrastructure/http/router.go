package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/triage", h.PersistTriageSession)
	rg.POST("/evaluate", h.EvaluateAutomaticDispatch)
	rg.POST("/recommendations/generate", h.GenerateRecommendations)
	rg.GET("/recommendations", h.ListRecommendations)
	rg.POST("/assignments", h.CreateAssignment)
	rg.GET("/assignments", h.ListAssignments)
	rg.GET("/assignments/:id", h.GetAssignmentByID)
	rg.PATCH("/assignments/:id/status", h.UpdateAssignmentStatus)
}
