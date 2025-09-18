#!/bin/bash

# AtCoder Discord Bot Setup Script

echo "🤖 AtCoder Discord Bot Setup"
echo "============================"
echo

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Check if .env file exists
if [ ! -f ".env" ]; then
    if [ -f ".env.example" ]; then
        echo "📁 Creating .env file from .env.example..."
        cp .env.example .env
        echo "✅ .env file created. Please edit it and add your Discord bot token."
        echo
        echo "To get a Discord bot token:"
        echo "1. Go to https://discord.com/developers/applications"
        echo "2. Create a new application"
        echo "3. Go to 'Bot' section and create a bot"
        echo "4. Copy the token and paste it in the .env file"
        echo "5. Enable 'Message Content Intent' in the bot settings"
        echo
    else
        echo "❌ .env.example file not found. Please create a .env file manually."
        exit 1
    fi
else
    echo "✅ .env file already exists"
fi

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod tidy
if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi

# Build the bot
echo "🔨 Building the bot..."
go build -o discord-bot .
if [ $? -eq 0 ]; then
    echo "✅ Bot built successfully"
else
    echo "❌ Failed to build the bot"
    exit 1
fi

echo
echo "🎉 Setup complete!"
echo
echo "Next steps:"
echo "1. Edit the .env file and add your Discord bot token"
echo "2. Run the bot with: ./discord-bot"
echo "   Or use: make run"
echo "   Or use Docker: make docker-compose-up"
echo
echo "Bot commands:"
echo "• !contest  - Show upcoming contests"
echo "• !next     - Show next contest"
echo "• !status   - Show bot status"
echo "• !help     - Show help message"
echo