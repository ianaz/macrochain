package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(ctx context.Context, redisHost string, redisPort int) (*RedisQueue, error) {
	slog.InfoContext(ctx, "Attempt to create new Redis queue", "host", redisHost, "port", redisPort)

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	queue := &RedisQueue{
		client: client,
	}

	slog.InfoContext(ctx, "Successfully created new Redis queue", "host", redisHost, "port", redisPort)
	return queue, nil
}

func (q *RedisQueue) Send(ctx context.Context, topic string, message Message) error {
	slog.InfoContext(ctx, "Attempt to send message", "topic", topic, "messageID", message.ID)

	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = q.client.Publish(ctx, topic, data).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	slog.InfoContext(ctx, "Successfully sent message", "topic", topic, "messageID", message.ID)
	return nil
}

func (q *RedisQueue) Subscribe(ctx context.Context, topic string) (<-chan Message, error) {
	slog.InfoContext(ctx, "Attempt to subscribe to topic", "topic", topic)

	// Create a subscription
	pubsub := q.client.Subscribe(ctx, topic)

	// Confirm that the subscription is working
	_, err := pubsub.Receive(ctx)
	if err != nil {
		pubsub.Close()
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	// Create message channel
	msgChan := make(chan Message, 100)

	// Create a done channel to signal when consumer is done
	done := make(chan struct{})

	// Start a goroutine to process messages
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.ErrorContext(context.Background(), "Panic in subscription goroutine",
					"topic", topic,
					"error", r,
				)
			}
			close(msgChan)
			slog.InfoContext(context.Background(), "Subscription closed", "topic", topic)
		}()

		channel := pubsub.Channel()

		for {
			select {
			case <-done:
				// Consumer has closed the channel, clean up
				err := pubsub.Unsubscribe(context.Background(), topic)
				if err != nil {
					slog.ErrorContext(context.Background(), "Failed to unsubscribe", "topic", topic, "error", err)
				}
				err = pubsub.Close()
				if err != nil {
					slog.ErrorContext(context.Background(), "Failed to close pubsub", "topic", topic, "error", err)
				}
				return

			case msg, ok := <-channel:
				if !ok {
					// Channel was closed, exit
					return
				}

				var message Message
				err := json.Unmarshal([]byte(msg.Payload), &message)
				if err != nil {
					slog.ErrorContext(context.Background(), "Failed to unmarshal message",
						"topic", topic,
						"error", err,
					)
					continue
				}

				// Log received message
				slog.InfoContext(context.Background(), "Received message from Redis",
					"topic", topic,
					"messageID", message.ID,
					"payload", string(message.Body),
				)

				// Try to send to the channel
				select {
				case msgChan <- message:
					// Message sent successfully
				case <-done:
					// Consumer has closed the channel, clean up
					return
				case <-time.After(1 * time.Second):
					slog.WarnContext(context.Background(), "Timed out sending message to consumer",
						"topic", topic,
					)
				}
			}
		}
	}()

	// Watch for consumer to close the channel
	go func() {
		// This goroutine will be leaked if the consumer doesn't actively close!
		// In a real system, we'd want more direct control or a different pattern.
		<-ctx.Done()
		close(done)
	}()

	slog.InfoContext(ctx, "Successfully subscribed to topic", "topic", topic)
	return msgChan, nil
}

func (q *RedisQueue) Unsubscribe(ctx context.Context, topic string) error {
	// Note: We now rely on context cancellation to clean up subscriptions
	// To unsubscribe, cancel the context that was used to create the subscription
	slog.InfoContext(ctx, "To unsubscribe: cancel the context used when subscribing", "topic", topic)
	return nil
}

func (q *RedisQueue) Ack(ctx context.Context, topic string, messageID string) error {
	slog.InfoContext(ctx, "Attempt to acknowledge message", "topic", topic, "messageID", messageID)
	slog.InfoContext(ctx, "Successfully acknowledged message", "topic", topic, "messageID", messageID)
	return nil
}

func (q *RedisQueue) Close() error {
	ctx := context.Background()
	slog.InfoContext(ctx, "Attempt to close Redis queue")

	// Close Redis client
	err := q.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close redis client: %w", err)
	}

	slog.InfoContext(ctx, "Successfully closed Redis queue")
	return nil
}
