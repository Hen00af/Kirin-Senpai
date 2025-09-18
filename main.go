package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Contest represents an AtCoder contest
type Contest struct {
	ID               string `json:"id"`
	StartEpochSecond int64  `json:"start_epoch_second"`
	DurationSecond   int64  `json:"duration_second"`
	Title            string `json:"title"`
	RateChange       string `json:"rate_change"`
}

// AtCoderAPIResponse represents the response from AtCoder API
type AtCoderAPIResponse []Contest

var (
	discord *discordgo.Session
	config  *Config
)

func main() {
	// Load configuration
	config = LoadConfig()

	if config.DiscordToken == "" {
		log.Fatal("DISCORD_TOKEN environment variable is required")
	}

	log.Printf("Starting AtCoder Contest Bot...")
	log.Printf("Command prefix: %s", config.CommandPrefix)
	log.Printf("Max contests to display: %d", config.MaxContests)

	// Create Discord session
	dg, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}
	discord = dg

	// Add message handler
	dg.AddHandler(messageCreate)

	// Indicate bot is ready
	dg.AddHandler(ready)

	// Open websocket connection
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}

	// Wait for interrupt signal
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "AtCoder contests | !contest")
	log.Printf("Bot is ready! Logged in as: %v#%v", event.User.Username, event.User.Discriminator)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Handle commands
	if strings.HasPrefix(m.Content, config.CommandPrefix+"contest") {
		handleContestCommand(s, m)
	} else if strings.HasPrefix(m.Content, config.CommandPrefix+"help") {
		handleHelpCommand(s, m)
	}
}

func handleContestCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Send typing indicator
	s.ChannelTyping(m.ChannelID)

	log.Printf("Contest command requested by user %s in channel %s", m.Author.Username, m.ChannelID)

	contests, err := getUpcomingContests()
	if err != nil {
		log.Printf("Error fetching contests: %v", err)
		s.ChannelMessageSend(m.ChannelID, "❌ Error fetching contest information. Please try again later.")
		return
	}

	if len(contests) == 0 {
		s.ChannelMessageSend(m.ChannelID, "📅 No upcoming contests found.")
		return
	}

	// Create embed with contest information
	embed := &discordgo.MessageEmbed{
		Title:       "🏆 Upcoming AtCoder Contests",
		Description: fmt.Sprintf("Here are the next %d upcoming AtCoder contests:", min(len(contests), config.MaxContests)),
		Color:       0x00ff00,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Data from AtCoder API • Updated every " + config.UpdateInterval.String(),
		},
	}

	// Add fields for each contest (limit to configured max)
	count := 0
	for _, contest := range contests {
		if count >= config.MaxContests {
			break
		}

		startTime := time.Unix(contest.StartEpochSecond, 0)
		duration := time.Duration(contest.DurationSecond) * time.Second

		// Calculate time until contest starts
		timeUntil := time.Until(startTime)
		var timeUntilStr string
		if timeUntil > 0 {
			timeUntilStr = fmt.Sprintf("\n**Starts in:** %s", formatDuration(timeUntil))
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: contest.Title,
			Value: fmt.Sprintf("**Start:** %s\n**Duration:** %s\n**Rate Change:** %s%s",
				startTime.Format("2006-01-02 15:04 MST"),
				formatDuration(duration),
				contest.RateChange,
				timeUntilStr),
			Inline: false,
		})
		count++
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	log.Printf("Sent contest information for %d contests", count)
}

func handleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🤖 AtCoder Contest Bot Help",
		Description: "I help you stay updated with AtCoder contests!",
		Color:       0x0099ff,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   config.CommandPrefix + "contest",
				Value:  "Show upcoming AtCoder contests",
				Inline: false,
			},
			{
				Name:   config.CommandPrefix + "help",
				Value:  "Show this help message",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "AtCoder Contest Bot • Data from kenkoooo.com/atcoder",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func getUpcomingContests() ([]Contest, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make HTTP request to AtCoder API
	resp, err := client.Get(config.AtCoderAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contests: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON response
	var allContests AtCoderAPIResponse
	if err := json.Unmarshal(body, &allContests); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Filter upcoming contests
	now := time.Now().Unix()
	var upcomingContests []Contest

	for _, contest := range allContests {
		if contest.StartEpochSecond > now {
			upcomingContests = append(upcomingContests, contest)
		}
	}

	// Sort by start time (earliest first)
	for i := 0; i < len(upcomingContests)-1; i++ {
		for j := i + 1; j < len(upcomingContests); j++ {
			if upcomingContests[i].StartEpochSecond > upcomingContests[j].StartEpochSecond {
				upcomingContests[i], upcomingContests[j] = upcomingContests[j], upcomingContests[i]
			}
		}
	}

	log.Printf("Found %d upcoming contests", len(upcomingContests))
	return upcomingContests, nil
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
