package infrastructure

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type DeviceTokenResolver interface {
	GetPushTokensByUserID(ctx context.Context, userID string) ([]string, error)
}

type Sender struct {
	fcm      *messaging.Client
	resolver DeviceTokenResolver
}

func NewSender(credentialsFile string, resolver DeviceTokenResolver) (*Sender, error) {
	app, err := firebase.NewApp(
		context.Background(),
		nil,
		option.WithAuthCredentialsFile(option.ServiceAccount, credentialsFile),
	)
	if err != nil {
		return nil, fmt.Errorf("init firebase app: %w", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("init firebase messaging client: %w", err)
	}

	return &Sender{
		fcm:      client,
		resolver: resolver,
	}, nil
}

func (s *Sender) SendSMS(ctx context.Context, to string, body string) error {
	// integrate SMS provider here
	return nil
}

func (s *Sender) SendEmail(ctx context.Context, to, subject, body string) error {
	// integrate email provider here
	return nil
}

func (s *Sender) SendPush(ctx context.Context, userID string, title, body string) error {
	if s == nil || s.fcm == nil {
		return errors.New("push sender not initialized")
	}
	if s.resolver == nil {
		return errors.New("device token resolver not configured")
	}
	if userID == "" {
		return errors.New("userID is required")
	}

	tokens, err := s.resolver.GetPushTokensByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("resolve device tokens: %w", err)
	}
	if len(tokens) == 0 {
		return nil
	}

	msg := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"user_id": userID,
			"type":    "GENERAL_NOTIFICATION",
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
		},
	}

	resp, err := s.fcm.SendEachForMulticast(ctx, msg)
	if err != nil {
		return fmt.Errorf("send push notification: %w", err)
	}

	if resp.FailureCount == len(tokens) {
		return fmt.Errorf("all push sends failed")
	}

	return nil
}
