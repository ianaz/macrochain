package queue

import (
	"context"
	"time"
)

type Message struct {
	ID        string
	Body      []byte
	Timestamp time.Time
	Metadata  map[string]string
}

type Queue interface {
	Send(ctx context.Context, topic string, message Message) error
	Subscribe(ctx context.Context, topic string) (<-chan Message, error)
	Unsubscribe(ctx context.Context, topic string) error
	Close() error
}
