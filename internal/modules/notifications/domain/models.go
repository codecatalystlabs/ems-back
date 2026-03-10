package domain

import "time"

type Notification struct {
	ID              string     `json:"id"`
	Type            string     `json:"type"`
	RecipientUserID *string    `json:"recipient_user_id,omitempty"`
	RecipientPhone  *string    `json:"recipient_phone,omitempty"`
	RecipientEmail  *string    `json:"recipient_email,omitempty"`
	Title           *string    `json:"title,omitempty"`
	Body            string     `json:"body"`
	Channel         string     `json:"channel"`
	LinkedEntityType *string   `json:"linked_entity_type,omitempty"`
	LinkedEntityID  *string    `json:"linked_entity_id,omitempty"`
	Status          string     `json:"status"`
	Attempts        int        `json:"attempts"`
	SentAt          *time.Time `json:"sent_at,omitempty"`
	ReadAt          *time.Time `json:"read_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

