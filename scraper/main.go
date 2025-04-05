package main

import (
	"context"
	"log/slog"
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

	// Main scraper loop
	for {
		// Example log for demonstration
		logger.InfoContext(ctx, "Scraper cycle starting")

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
