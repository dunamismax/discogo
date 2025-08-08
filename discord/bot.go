// Package discord provides Discord bot functionality.
package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dunamismax/discogo/config"
	"github.com/dunamismax/discogo/errors"
	"github.com/dunamismax/discogo/logging"
	"github.com/dunamismax/discogo/metrics"
)

// Bot represents a Discord bot instance with all necessary components.
type Bot struct {
	session         *discordgo.Session
	config          *config.Config
	commandHandlers map[string]CommandHandler
}

// CommandHandler represents a function that handles Discord bot commands.
type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error

// NewBot creates a new Discord bot instance.
func NewBot(cfg *config.Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, errors.NewDiscordError("failed to create Discord session", err)
	}

	bot := &Bot{
		session:         session,
		config:          cfg,
		commandHandlers: make(map[string]CommandHandler),
	}

	// Register command handlers.
	bot.registerCommands()

	// Add message handler.
	session.AddHandler(bot.messageCreate)

	// Set intents.
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	return bot, nil
}

// Start starts the Discord bot.
func (b *Bot) Start() error {
	logger := logging.WithComponent("discord")
	logger.Info("Starting bot", "bot_name", b.config.BotName)

	err := b.session.Open()
	if err != nil {
		return errors.NewDiscordError("failed to open Discord session", err)
	}

	logger.Info("Bot is now running", "username", b.session.State.User.Username)

	return nil
}

// Stop stops the Discord bot.
func (b *Bot) Stop() error {
	logger := logging.WithComponent("discord")
	logger.Info("Stopping bot", "bot_name", b.config.BotName)

	if err := b.session.Close(); err != nil {
		return errors.NewDiscordError("failed to close Discord session", err)
	}

	return nil
}

// registerCommands registers all command handlers.
func (b *Bot) registerCommands() {
	b.commandHandlers["ping"] = b.handlePing
	b.commandHandlers["help"] = b.handleHelp
	b.commandHandlers["stats"] = b.handleStats
}

// messageCreate handles incoming messages.
func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from bots.
	if m.Author.Bot {
		return
	}

	// Check if message starts with command prefix.
	if !strings.HasPrefix(m.Content, b.config.CommandPrefix) {
		return
	}

	// Remove prefix and split into command and args.
	content := strings.TrimPrefix(m.Content, b.config.CommandPrefix)

	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]

	// Handle specific commands.
	if handler, exists := b.commandHandlers[command]; exists {
		if err := handler(s, m, args); err != nil {
			logger := logging.WithComponent("discord").With(
				"user_id", m.Author.ID,
				"username", m.Author.Username,
				"command", command,
			)
			logging.LogError(logger, err, "Command execution failed")
			metrics.RecordCommand(false)
			metrics.RecordError(err)
			b.sendErrorMessage(s, m.ChannelID, "Sorry, something went wrong processing your command.")
		} else {
			metrics.RecordCommand(true)
			logging.LogDiscordCommand(m.Author.ID, m.Author.Username, command, true)
		}

		return
	}

	// If no specific handler found, send unknown command message.
	b.sendErrorMessage(s, m.ChannelID, fmt.Sprintf("Unknown command: %s%s. Use %shelp for available commands.", b.config.CommandPrefix, command, b.config.CommandPrefix))
}

// handlePing handles the !ping command.
func (b *Bot) handlePing(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	logger := logging.WithComponent("discord").With(
		"user_id", m.Author.ID,
		"username", m.Author.Username,
		"command", "ping",
	)
	logger.Info("Handling ping command")

	embed := &discordgo.MessageEmbed{
		Title:       "Pong! 🏓",
		Description: "Bot is online and responding!",
		Color:       0x00FF00, // Green color
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return errors.NewDiscordError("failed to send ping response", err)
	}

	return nil
}

// sendErrorMessage sends an error message to a Discord channel.
func (b *Bot) sendErrorMessage(s *discordgo.Session, channelID, message string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Error",
		Description: message,
		Color:       0xE74C3C, // Red color.
	}

	if _, err := s.ChannelMessageSendEmbed(channelID, embed); err != nil {
		logger := logging.WithComponent("discord")
		logger.Error("Failed to send error message", "error", err)
	}
}

// handleHelp handles the !help command.
func (b *Bot) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	logger := logging.WithComponent("discord").With(
		"user_id", m.Author.ID,
		"username", m.Author.Username,
		"command", "help",
	)
	logger.Info("Showing help information")

	embed := &discordgo.MessageEmbed{
		Title:       "Discord Bot Help",
		Description: "A generic Discord bot template built with Go!",
		Color:       0x3498DB, // Blue color.
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("%sping", b.config.CommandPrefix),
				Value:  "Check if the bot is online and responding",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%shelp", b.config.CommandPrefix),
				Value:  "Show this help message",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%sstats", b.config.CommandPrefix),
				Value:  "Show bot performance statistics",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "🚀 Built with Go, DiscordGo, and Mage - Ready for customization!",
		},
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return errors.NewDiscordError("failed to send help message", err)
	}

	return nil
}

// handleStats handles the !stats command.
func (b *Bot) handleStats(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	logger := logging.WithComponent("discord").With(
		"user_id", m.Author.ID,
		"username", m.Author.Username,
		"command", "stats",
	)
	logger.Info("Showing bot statistics")

	summary := metrics.Get().GetSummary()
	uptime := time.Duration(summary.UptimeSeconds * float64(time.Second))

	// Format uptime nicely.
	uptimeStr := formatDuration(uptime)

	embed := &discordgo.MessageEmbed{
		Title: "Bot Statistics",
		Color: 0x2ECC71, // Green color.
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "📊 Commands",
				Value: fmt.Sprintf("Total: %d\nSuccessful: %d\nFailed: %d\nSuccess Rate: %.1f%%",
					summary.CommandsTotal, summary.CommandsSuccessful, summary.CommandsFailed, summary.CommandSuccessRate),
				Inline: true,
			},
			{
				Name: "🌐 API Requests",
				Value: fmt.Sprintf("Total: %d\nSuccess Rate: %.1f%%\nAvg Response: %.0fms",
					summary.APIRequestsTotal, summary.APISuccessRate, summary.AverageResponseTime),
				Inline: true,
			},
			{
				Name: "⚡ Performance",
				Value: fmt.Sprintf("Commands/sec: %.2f\nAPI Requests/sec: %.2f",
					summary.CommandsPerSecond, summary.APIRequestsPerSecond),
				Inline: true,
			},
			{
				Name:   "⏱️ Uptime",
				Value:  uptimeStr,
				Inline: true,
			},
			{
				Name:   "🚀 Started",
				Value:  fmt.Sprintf("<t:%d:R>", time.Now().Add(-uptime).Unix()),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Statistics since bot startup",
		},
	}

	// Add error information if there are errors.
	if len(summary.ErrorsByType) > 0 {
		errorInfo := make([]string, 0, len(summary.ErrorsByType))
		for errorType, count := range summary.ErrorsByType {
			if count > 0 {
				errorInfo = append(errorInfo, fmt.Sprintf("%s: %d", string(errorType), count))
			}
		}

		if len(errorInfo) > 0 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "⚠️ Errors",
				Value:  strings.Join(errorInfo, "\n"),
				Inline: false,
			})
		}
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return errors.NewDiscordError("failed to send stats message", err)
	}

	return nil
}

// formatDuration formats a duration into a human-readable string.
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	case hours > 0:
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	case minutes > 0:
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	default:
		return fmt.Sprintf("%ds", seconds)
	}
}
