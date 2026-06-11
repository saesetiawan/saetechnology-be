package broker

import "context"

type Message struct {
	Key   string
	Value []byte
}

type Handler func(ctx context.Context, msg Message) error

type Consumer interface {
	Start(ctx context.Context, topic string, handler Handler) error
}
