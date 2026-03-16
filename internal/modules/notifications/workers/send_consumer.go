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

type Sender interface {
	SendSMS(ctx context.Context, to string, body string) error
	SendEmail(ctx context.Context, to, subject, body string) error
	SendPush(ctx context.Context, userID string, title, body string) error
}

type SendConsumer struct {
	repo   notificationsapp.Repository
	sender Sender
	bus    events.Publisher
	log    *zap.Logger
}

func NewSendConsumer(repo notificationsapp.Repository, sender Sender, bus events.Publisher, log *zap.Logger) *SendConsumer {
	return &SendConsumer{repo: repo, sender: sender, bus: bus, log: log}
}

func (c *SendConsumer) Handle(ctx context.Context, evt events.Event) error {
	notificationID, _ := evt.Payload["notification_id"].(string)
	channel, _ := evt.Payload["channel"].(string)
	recipientPhone, _ := evt.Payload["recipient_phone"].(string)
	recipientEmail, _ := evt.Payload["recipient_email"].(string)
	title, _ := evt.Payload["title"].(string)
	body, _ := evt.Payload["body"].(string)
	recipientUserID, _ := evt.Payload["recipient_user_id"].(string)

	_ = c.repo.IncrementAttempts(ctx, notificationID)

	var err error
	switch channel {
	case "SMS":
		err = c.sender.SendSMS(ctx, recipientPhone, body)
	case "EMAIL":
		err = c.sender.SendEmail(ctx, recipientEmail, title, body)
	case "PUSH":
		err = c.sender.SendPush(ctx, recipientUserID, title, body)
	case "IN_APP":
		err = nil
	default:
		err = nil
	}

	if err != nil {
		_ = c.repo.MarkFailed(ctx, notificationID)
		_ = c.bus.Publish(ctx, constants.TopicNotificationFailed, events.Event{
			ID:          uuid.NewString(),
			Topic:       constants.TopicNotificationFailed,
			AggregateID: notificationID,
			Type:        constants.TopicNotificationFailed,
			OccurredAt:  time.Now().UTC(),
			Payload: map[string]any{
				"notification_id": notificationID,
				"error":           err.Error(),
			},
		})
		return err
	}

	_ = c.repo.MarkSent(ctx, notificationID)
	_ = c.bus.Publish(ctx, constants.TopicNotificationSent, events.Event{
		ID:          uuid.NewString(),
		Topic:       constants.TopicNotificationSent,
		AggregateID: notificationID,
		Type:        constants.TopicNotificationSent,
		OccurredAt:  time.Now().UTC(),
		Payload: map[string]any{
			"notification_id": notificationID,
		},
	})
	return nil
}
