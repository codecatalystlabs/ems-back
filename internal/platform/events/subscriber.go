package events

import "context"

type Handler interface {
	Handle(ctx context.Context, event Event) error
}
