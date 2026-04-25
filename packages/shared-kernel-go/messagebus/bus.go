package messagebus

import "context"

type Handler func(ctx context.Context, event Event) error

type MessageBus interface {
	Publish(ctx context.Context, topic string, event Event) error
	Subscribe(queue string, handler Handler) error
}
