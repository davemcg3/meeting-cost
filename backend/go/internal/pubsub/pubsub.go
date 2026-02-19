package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// PubSub handles publishing and subscribing to events via Redis.
type PubSub interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string) <-chan string
}

type redisPubSub struct {
	client *redis.Client
}

// NewRedisPubSub creates a new Redis-based PubSub.
func NewRedisPubSub(client *redis.Client) PubSub {
	return &redisPubSub{
		client: client,
	}
}

func (p *redisPubSub) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}

	return p.client.Publish(ctx, channel, data).Err()
}

func (p *redisPubSub) Subscribe(ctx context.Context, channel string) <-chan string {
	ch := make(chan string)
	pubsub := p.client.Subscribe(ctx, channel)

	go func() {
		defer pubsub.Close()
		defer close(ch)

		for msg := range pubsub.Channel() {
			ch <- msg.Payload
		}
	}()

	return ch
}
