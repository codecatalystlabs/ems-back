package infrastructure

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("", h.Register)
	rg.GET("", h.List)
	rg.GET("/:id", h.GetByID)
	rg.PATCH("/:id", h.Update)
	rg.DELETE("/:id", h.Deactivate)
}
