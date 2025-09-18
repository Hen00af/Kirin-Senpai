package main

import (
	"testing"
	"time"
)

func TestGetUpcomingContests(t *testing.T) {
	// Initialize config for testing
	config = &Config{
		AtCoderAPIURL: "https://kenkoooo.com/atcoder/resources/contests.json",
		MaxContests:   5,
	}

	contests, err := getUpcomingContests()
	if err != nil {
		t.Fatalf("Failed to get upcoming contests: %v", err)
	}

	// We should get a list (could be empty if no upcoming contests)
	if contests == nil {
		t.Fatal("Expected contests slice, got nil")
	}

	// If we have contests, verify they have required fields
	for i, contest := range contests {
		if contest.ID == "" {
			t.Errorf("Contest %d has empty ID", i)
		}
		if contest.Title == "" {
			t.Errorf("Contest %d has empty title", i)
		}
		if contest.StartEpochSecond <= 0 {
			t.Errorf("Contest %d has invalid start time: %d", i, contest.StartEpochSecond)
		}
		if contest.DurationSecond <= 0 {
			t.Errorf("Contest %d has invalid duration: %d", i, contest.DurationSecond)
		}

		// Verify start time is in the future
		if contest.StartEpochSecond <= time.Now().Unix() {
			t.Errorf("Contest %d start time is not in the future", i)
		}
	}

	// If we have multiple contests, verify they are sorted by start time
	for i := 1; i < len(contests); i++ {
		if contests[i-1].StartEpochSecond > contests[i].StartEpochSecond {
			t.Errorf("Contests are not sorted by start time")
			break
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Minute, "30m"},
		{1 * time.Hour, "1h 0m"},
		{1*time.Hour + 30*time.Minute, "1h 30m"},
		{2*time.Hour + 45*time.Minute, "2h 45m"},
		{5 * time.Minute, "5m"},
		{25 * time.Hour, "1d 1h 0m"},
		{48 * time.Hour, "2d 0h 0m"},
		{-30 * time.Minute, "30m"}, // Test negative duration
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %s, expected %s", test.duration, result, test.expected)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	// Test default config loading
	config := LoadConfig()
	
	// Check defaults
	if config.AtCoderAPIURL != "https://kenkoooo.com/atcoder/resources/contests.json" {
		t.Errorf("Expected default AtCoder API URL, got %s", config.AtCoderAPIURL)
	}
	
	if config.MaxContests != 5 {
		t.Errorf("Expected default max contests 5, got %d", config.MaxContests)
	}
	
	if config.CommandPrefix != "!" {
		t.Errorf("Expected default command prefix '!', got %s", config.CommandPrefix)
	}
	
	if config.UpdateInterval != 10*time.Minute {
		t.Errorf("Expected default update interval 10m, got %v", config.UpdateInterval)
	}
}

func TestMinFunction(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{5, 3, 3},
		{1, 10, 1},
		{7, 7, 7},
		{0, 5, 0},
		{-1, 3, -1},
	}
	
	for _, test := range tests {
		result := min(test.a, test.b)
		if result != test.expected {
			t.Errorf("min(%d, %d) = %d, expected %d", test.a, test.b, result, test.expected)
		}
	}
}
