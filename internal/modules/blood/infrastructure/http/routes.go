package http

import (
	rbacmiddleware "dispatch/internal/modules/rbac/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/requisitions", rbacmiddleware.RequirePermission(nil, "incidents.create"), h.RaiseRequisition)
	rg.GET("/requisitions", rbacmiddleware.RequirePermission(nil, "incidents.read"), h.ListRequisitions)
	rg.POST("/requisitions/:id/broadcast", rbacmiddleware.RequirePermission(nil, "dispatch.assign"), h.Broadcast)
	rg.GET("/requisitions/:id/offers", rbacmiddleware.RequirePermission(nil, "incidents.read"), h.ListOffers)
	rg.POST("/offers", rbacmiddleware.RequirePermission(nil, "incidents.create"), h.CreateOffer)
	rg.POST("/requisitions/:id/offers/:offerId/accept", rbacmiddleware.RequirePermission(nil, "dispatch.assign"), h.AcceptOffer)
	rg.POST("/pickup-assignments", rbacmiddleware.RequirePermission(nil, "dispatch.assign"), h.AssignPickup)
	rg.POST("/pickup-assignments/:assignmentId/collect", rbacmiddleware.RequirePermission(nil, "dispatch.update_status"), h.MarkCollected)
	rg.POST("/pickup-assignments/:assignmentId/deliver", rbacmiddleware.RequirePermission(nil, "dispatch.update_status"), h.MarkDelivered)
}
