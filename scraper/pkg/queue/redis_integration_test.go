//go:build integration
// +build integration

package queue

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestRedisQueueIntegration(t *testing.T) {
	// Get Redis connection info from environment variables or use defaults
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPortStr := getEnv("REDIS_PORT", "6379")
	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		t.Fatalf("Invalid Redis port: %v", err)
	}

	ctx := context.Background()

	// Create Redis queue
	queue, err := NewRedisQueue(ctx, redisHost, redisPort)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()

	// Test topic
	topic := "test-topic-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// Subscribe to test topic
	messages, err := queue.Subscribe(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to subscribe to topic: %v", err)
	}

	// Give Redis some time to set up the subscription
	time.Sleep(500 * time.Millisecond)

	// Test message
	testMessage := Message{
		Body:     []byte("test message"),
		Metadata: map[string]string{"test": "true"},
	}

	// Send message
	err = queue.Send(ctx, topic, testMessage)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Wait for message to be received
	select {
	case receivedMsg := <-messages:
		// Verify message content
		if string(receivedMsg.Body) != string(testMessage.Body) {
			t.Errorf("Expected message body %q, got %q", testMessage.Body, receivedMsg.Body)
		}
		if receivedMsg.Metadata["test"] != testMessage.Metadata["test"] {
			t.Errorf("Expected metadata %v, got %v", testMessage.Metadata, receivedMsg.Metadata)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for message")
	}

	// Test unsubscribe
	err = queue.Unsubscribe(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}
}

func TestMultipleSubscribersIntegration(t *testing.T) {
	// Get Redis connection info from environment variables or use defaults
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPortStr := getEnv("REDIS_PORT", "6379")
	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		t.Fatalf("Invalid Redis port: %v", err)
	}

	ctx := context.Background()

	// Create Redis queue
	queue, err := NewRedisQueue(ctx, redisHost, redisPort)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()

	// Test topic
	topic := "test-topic-multi-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// Wait for Redis to be fully ready for subscribing
	time.Sleep(200 * time.Millisecond)

	// Create multiple subscribers
	messages1, err := queue.Subscribe(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to create first subscriber: %v", err)
	}

	// Wait between subscriptions to ensure proper setup
	time.Sleep(200 * time.Millisecond)

	messages2, err := queue.Subscribe(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to create second subscriber: %v", err)
	}

	// Wait for subscriptions to be fully set up before sending
	time.Sleep(500 * time.Millisecond)

	// Test message
	testMessage := Message{
		Body:     []byte("multi-subscriber test"),
		Metadata: map[string]string{"multi": "true"},
	}

	// Send message
	err = queue.Send(ctx, topic, testMessage)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Check that both subscribers receive the message
	// Subscriber 1
	select {
	case receivedMsg := <-messages1:
		if string(receivedMsg.Body) != string(testMessage.Body) {
			t.Errorf("Subscriber 1: Expected message body %q, got %q", testMessage.Body, receivedMsg.Body)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Subscriber 1: Timed out waiting for message")
	}

	// Subscriber 2
	select {
	case receivedMsg := <-messages2:
		if string(receivedMsg.Body) != string(testMessage.Body) {
			t.Errorf("Subscriber 2: Expected message body %q, got %q", testMessage.Body, receivedMsg.Body)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Subscriber 2: Timed out waiting for message")
	}

	// Cleanup
	err = queue.Unsubscribe(ctx, topic)
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
