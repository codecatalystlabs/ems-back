package application

type CreateNotificationRequest struct {
	Type             string  `json:"type" binding:"required"`
	RecipientUserID  *string `json:"recipient_user_id,omitempty"`
	RecipientPhone   *string `json:"recipient_phone,omitempty"`
	RecipientEmail   *string `json:"recipient_email,omitempty"`
	Title            *string `json:"title,omitempty"`
	Body             string  `json:"body" binding:"required"`
	Channel          string  `json:"channel" binding:"required"`
	LinkedEntityType *string `json:"linked_entity_type,omitempty"`
	LinkedEntityID   *string `json:"linked_entity_id,omitempty"`
}

type UpdateNotificationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
