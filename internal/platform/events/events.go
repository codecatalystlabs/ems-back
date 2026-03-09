package events

import "time"

type Event struct {
	ID          string         `json:"id"`
	Topic       string         `json:"topic"`
	AggregateID string         `json:"aggregate_id"`
	Type        string         `json:"type"`
	OccurredAt  time.Time      `json:"occurred_at"`
	Payload     map[string]any `json:"payload"`
}
