package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/requisitions", h.RaiseRequisition)
	rg.GET("/requisitions", h.ListRequisitions)
	rg.POST("/requisitions/:id/broadcast", h.Broadcast)
	rg.GET("/requisitions/:id/offers", h.ListOffers)
	rg.POST("/offers", h.CreateOffer)
	rg.POST("/requisitions/:id/offers/:offerId/accept", h.AcceptOffer)
	rg.POST("/pickup-assignments", h.AssignPickup)
	rg.POST("/pickup-assignments/:assignmentId/collect", h.MarkCollected)
	rg.POST("/pickup-assignments/:assignmentId/deliver", h.MarkDelivered)
}
