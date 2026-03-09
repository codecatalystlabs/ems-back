package events

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, event Event) error
}
