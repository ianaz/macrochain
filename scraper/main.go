package main

import (
	"context"
	"log/slog"
	"macrochain/scraper/pkg/queue"
	"time"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	logger := SetupLogger(config.LogLevel)
	slog.SetDefault(logger)

	ctx := context.Background()
	logger.InfoContext(ctx, "Starting Macrochain scraper",
		"db_host", config.DBHost,
		"redis_host", config.RedisHost,
		"scrape_interval", config.ScrapeInterval)

	redisQueue, err := queue.NewRedisQueue(ctx, config.RedisHost, config.RedisPort)
	if err != nil {
		panic("Failed to connect to Redis queue: " + err.Error())
	}
	defer redisQueue.Close()

	// Main scraper loop
	for {
		// Example log for demonstration
		logger.InfoContext(ctx, "Scraper cycle starting")

		// Example of sending a message to a queue
		message := queue.Message{
			Body:     []byte("Scraper cycle started"),
			Metadata: map[string]string{"source": "scraper", "type": "cycle_start"},
		}

		err := redisQueue.Send(ctx, "scraper_events", message)
		if err != nil {
			logger.ErrorContext(ctx, "Failed to send message to queue", "error", err)
		}

		// TODO: Implement scrapers for different data sources
		// - FED data
		// - SNB data
		// - Ethereum on-chain data
		// - DeFi protocols

		logger.InfoContext(ctx, "Scraper cycle completed")

		// Sleep until next cycle
		time.Sleep(time.Duration(config.ScrapeInterval) * time.Second)
	}
}
