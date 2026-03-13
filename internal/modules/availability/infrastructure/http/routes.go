package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	shift := rg.Group("/shifts")
	shift.POST("", h.CreateShift)
	shift.GET("", h.ListShifts)
	shift.GET("/:id", h.GetShiftByID)
	shift.PATCH("/:id", h.UpdateShift)

	avail := rg.Group("/users")
	avail.GET("/availability", h.ListAvailability)
	avail.POST("/availability", h.UpsertAvailability)
	avail.GET("/:userId/availability", h.GetAvailabilityByUserID)
	avail.POST("/presence-logs", h.CreatePresenceLog)
	avail.GET("/presence-logs", h.ListPresenceLogs)
}
