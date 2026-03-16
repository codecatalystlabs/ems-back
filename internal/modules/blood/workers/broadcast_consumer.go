package workers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	notificationsapp "dispatch/internal/modules/notifications/application"
	"dispatch/internal/platform/events"
	"dispatch/internal/shared/constants"
)

type BroadcastTargetFinder interface {
	FindBroadcastRecipients(ctx context.Context, bloodRequisitionID string) ([]BroadcastRecipient, error)
}

type BroadcastRecipient struct {
	UserID *string
	Phone  string
	Email  string
	Name   string
}

type BroadcastConsumer struct {
	notifier *notificationsapp.Service
	finder   BroadcastTargetFinder
	bus      events.Publisher
	log      *zap.Logger
}

func NewBroadcastConsumer(
	notifier *notificationsapp.Service,
	finder BroadcastTargetFinder,
	bus events.Publisher,
	log *zap.Logger,
) *BroadcastConsumer {
	return &BroadcastConsumer{
		notifier: notifier,
		finder:   finder,
		bus:      bus,
		log:      log,
	}
}

func (c *BroadcastConsumer) Handle(ctx context.Context, evt events.Event) error {
	bloodReqID, _ := evt.Payload["blood_requisition_id"].(string)
	requestMessage, _ := evt.Payload["request_message"].(string)

	_ = c.bus.Publish(ctx, constants.TopicBloodBroadcastDispatching, events.Event{
		ID:          uuid.NewString(),
		Topic:       constants.TopicBloodBroadcastDispatching,
		AggregateID: bloodReqID,
		Type:        constants.TopicBloodBroadcastDispatching,
		OccurredAt:  time.Now().UTC(),
		Payload:     evt.Payload,
	})

	recipients, err := c.finder.FindBroadcastRecipients(ctx, bloodReqID)
	if err != nil {
		_ = c.bus.Publish(ctx, constants.TopicBloodBroadcastFailed, events.Event{
			ID:          uuid.NewString(),
			Topic:       constants.TopicBloodBroadcastFailed,
			AggregateID: bloodReqID,
			Type:        constants.TopicBloodBroadcastFailed,
			OccurredAt:  time.Now().UTC(),
			Payload: map[string]any{
				"blood_requisition_id": bloodReqID,
				"error":                err.Error(),
			},
		})
		return err
	}

	sentCount := 0
	for _, r := range recipients {
		title := "Urgent Blood Request"

		_, nerr := c.notifier.Create(
			ctx,
			"BLOOD_REQUEST_BROADCAST",
			"SMS",
			r.UserID,
			&r.Phone,
			&r.Email,
			&title,
			&requestMessage,
			"blood_requisition",
			&bloodReqID,
		)
		if nerr != nil {
			c.log.Error("queue blood broadcast notification", zap.Error(nerr), zap.String("blood_requisition_id", bloodReqID))
			continue
		}
		sentCount++
	}

	_ = c.bus.Publish(ctx, constants.TopicBloodBroadcastSent, events.Event{
		ID:          uuid.NewString(),
		Topic:       constants.TopicBloodBroadcastSent,
		AggregateID: bloodReqID,
		Type:        constants.TopicBloodBroadcastSent,
		OccurredAt:  time.Now().UTC(),
		Payload: map[string]any{
			"blood_requisition_id": bloodReqID,
			"recipient_count":      sentCount,
		},
	})

	return nil
}
