package main

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for the scraper
type Config struct {
	LogLevel       string `mapstructure:"LOG_LEVEL"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         int    `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	RedisHost      string `mapstructure:"REDIS_HOST"`
	RedisPort      int    `mapstructure:"REDIS_PORT"`
	ScrapeInterval int    `mapstructure:"SCRAPE_INTERVAL"`
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "macrochain")
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("SCRAPE_INTERVAL", 3600) // 1 hour in seconds

	// Read from environment variables
	v.AutomaticEnv()

	var config Config
	err := v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
