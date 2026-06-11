package broker

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, msg Message) error
}
