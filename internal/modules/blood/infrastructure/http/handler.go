package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	bloodapp "dispatch/internal/modules/blood/application"
	platformdb "dispatch/internal/platform/db"
	"dispatch/internal/platform/httpx"
)

type Handler struct{ service *bloodapp.Service }

func NewHandler(service *bloodapp.Service) *Handler { return &Handler{service: service} }

func (h *Handler) RaiseRequisition(c *gin.Context) {
	var req bloodapp.CreateBloodRequisitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.RaiseRequisition(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Created(c, out)
}

func (h *Handler) Broadcast(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		DestinationLat *float64 `json:"destination_lat"`
		DestinationLon *float64 `json:"destination_lon"`
	}
	_ = c.ShouldBindJSON(&payload)
	out, err := h.service.BroadcastRequisition(c.Request.Context(), id, payload.DestinationLat, payload.DestinationLon)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

func (h *Handler) ListRequisitions(c *gin.Context) {
	p := platformdb.ParsePagination(c.Request.URL.Query(), map[string]string{
		"created_at":      "br.created_at",
		"status":          "br.status",
		"urgency_level":   "br.urgency_level",
		"units_requested": "br.units_requested",
	}, map[string]struct{}{
		"status":        {},
		"urgency_level": {},
	})
	out, err := h.service.ListRequisitions(c.Request.Context(), p)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

func (h *Handler) CreateOffer(c *gin.Context) {
	var req bloodapp.CreateBloodOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.CreateOffer(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Created(c, out)
}

func (h *Handler) ListOffers(c *gin.Context) {
	requisitionID := c.Param("id")
	p := platformdb.ParsePagination(c.Request.URL.Query(), map[string]string{
		"created_at":    "bro.created_at",
		"status":        "bro.status",
		"units_offered": "bro.units_offered",
	}, map[string]struct{}{
		"status": {},
	})
	out, err := h.service.ListOffers(c.Request.Context(), requisitionID, p)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, out)
}

func (h *Handler) AcceptOffer(c *gin.Context) {
	requisitionID := c.Param("id")
	offerID := c.Param("offerId")
	var body struct {
		ActorUserID *string `json:"actor_user_id"`
	}
	_ = c.ShouldBindJSON(&body)
	if err := h.service.AcceptOffer(c.Request.Context(), requisitionID, offerID, body.ActorUserID); err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, gin.H{"message": "offer accepted"})
}

func (h *Handler) AssignPickup(c *gin.Context) {
	var req bloodapp.AssignBloodPickupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	out, err := h.service.AssignPickup(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.Created(c, out)
}

func (h *Handler) MarkCollected(c *gin.Context) {
	assignmentID := c.Param("assignmentId")
	var body struct {
		BloodRequisitionID string  `json:"blood_requisition_id" binding:"required"`
		ActorUserID        *string `json:"actor_user_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.MarkCollected(c.Request.Context(), assignmentID, body.BloodRequisitionID, body.ActorUserID); err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, gin.H{"message": "blood marked collected"})
}

func (h *Handler) MarkDelivered(c *gin.Context) {
	assignmentID := c.Param("assignmentId")
	var body struct {
		BloodRequisitionID string  `json:"blood_requisition_id" binding:"required"`
		ActorUserID        *string `json:"actor_user_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.MarkDelivered(c.Request.Context(), assignmentID, body.BloodRequisitionID, body.ActorUserID); err != nil {
		httpx.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.OK(c, gin.H{"message": "blood marked delivered"})
}
