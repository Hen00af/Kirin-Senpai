package main

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	DiscordToken   string
	AtCoderAPIURL  string
	UpdateInterval time.Duration
	MaxContests    int
	CommandPrefix  string
	LogLevel       string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	config := &Config{
		DiscordToken:   getEnv("DISCORD_TOKEN", ""),
		AtCoderAPIURL:  getEnv("ATCODER_API_URL", "https://kenkoooo.com/atcoder/resources/contests.json"),
		UpdateInterval: getEnvDuration("UPDATE_INTERVAL", 10*time.Minute),
		MaxContests:    getEnvInt("MAX_CONTESTS", 5),
		CommandPrefix:  getEnv("COMMAND_PREFIX", "!"),
		LogLevel:       getEnv("LOG_LEVEL", "INFO"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
