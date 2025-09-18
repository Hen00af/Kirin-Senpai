# AtCoder Contest Discord Bot

A Discord bot written in Go that fetches and displays upcoming AtCoder contest information using the AtCoder API.

## Features

- 🏆 Fetch upcoming AtCoder contests
- 📅 Display contest start times and durations
- 🎯 Show rate change information
- 🤖 Simple Discord slash commands
- 🐳 Docker support for easy deployment

## Commands

- `!contest` - Display upcoming AtCoder contests
- `!help` - Show available commands

## Setup

### Prerequisites

- Go 1.21 or higher
- Discord Bot Token

### Getting Discord Bot Token

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to "Bot" section
4. Create a bot and copy the token
5. Under "Privileged Gateway Intents", enable "Message Content Intent"
6. Invite the bot to your server with appropriate permissions

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd discord-atcoder-bot
```

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Edit `.env` and add your Discord bot token:
```bash
DISCORD_TOKEN=your_discord_bot_token_here
```

4. Run the bot:
```bash
# Install dependencies
go mod tidy

# Run the bot
go run main.go
```

### Using Docker

1. Build the Docker image:
```bash
docker build -t discord-atcoder-bot .
```

2. Run the container:
```bash
docker run -e DISCORD_TOKEN=your_discord_bot_token_here discord-atcoder-bot
```

Or using docker-compose:

```yaml
version: '3.8'
services:
  discord-bot:
    build: .
    environment:
      - DISCORD_TOKEN=your_discord_bot_token_here
    restart: unless-stopped
```

## API Reference

This bot uses the AtCoder API (via kenkoooo.com) to fetch contest information:
- Endpoint: `https://kenkoooo.com/atcoder/resources/contests.json`
- Returns: Array of contest objects with ID, start time, duration, title, and rate change info

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the bot functionality
5. Submit a pull request

## License

This project is open source and available under the MIT License.