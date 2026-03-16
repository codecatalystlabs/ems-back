package infrastructure

import "context"

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) SendSMS(ctx context.Context, to string, body string) error {
	// integrate Africa's Talking, Twilio, Infobip, etc.
	return nil
}

func (s *Sender) SendEmail(ctx context.Context, to, subject, body string) error {
	// integrate SMTP/SES/SendGrid
	return nil
}

func (s *Sender) SendPush(ctx context.Context, userID string, title, body string) error {
	// integrate FCM/APNS
	return nil
}
