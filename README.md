<p align="center">
  <img src="https://github.com/dunamismax/images/blob/main/golang/discord-bots/generic-bot.png" alt="Discord Bot" width="300" />
</p>

<p align="center">
  <a href="https://github.com/dunamismax/discogo">
    <img src="https://readme-typing-svg.demolab.com/?font=Fira+Code&size=24&pause=1000&color=00ADD8&center=true&vCenter=true&width=900&lines=Discord+Bot+Template+in+Go;Clean+Architecture+with+DiscordGo;Structured+Logging+and+Metrics;Environment+Configuration+Management;Auto-Restart+Development+with+Mage;Single+Binary+Deployments;Production+Ready+Template" alt="Typing SVG" />
  </a>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.24+-00ADD8.svg?logo=go" alt="Go Version"></a>
  <a href="https://github.com/bwmarrin/discordgo"><img src="https://img.shields.io/badge/Discord-DiscordGo-5865F2.svg?logo=discord&logoColor=white" alt="DiscordGo"></a>
  <a href="https://magefile.org/"><img src="https://img.shields.io/badge/Build-Mage-purple.svg?logo=go" alt="Mage"></a>
  <a href="https://pkg.go.dev/log/slog"><img src="https://img.shields.io/badge/Logging-slog-00ADD8.svg?logo=go" alt="Go slog"></a>
  <a href="https://github.com/spf13/viper"><img src="https://img.shields.io/badge/Config-Environment-00ADD8.svg?logo=go" alt="Environment Config"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License"></a>
</p>

---

## About

A clean, production-ready Discord bot template built in Go. Perfect starting point for creating Discord bots with proper architecture, logging, metrics, and development tools.

**Features:**

* **Clean Architecture** – Well-organized packages with clear separation of concerns
* **Rich Discord Integration** – Built on DiscordGo with proper command handling
* **Structured Logging** – Using Go's native slog with configurable levels
* **Performance Metrics** – Built-in command and error tracking
* **Development Tools** – Auto-restart, build scripts, and quality checks with Mage
* **Easy Configuration** – Environment variables with validation and defaults
* **Production Ready** – Single-binary builds with graceful shutdown

---

## Quick Start

```bash
git clone https://github.com/dunamismax/discogo.git
cd discogo
go mod tidy
go install github.com/magefile/mage@latest
cp .env.example .env  # Add your Discord bot token
mage setup
mage dev
```

**Requirements:** Go 1.24+, Discord Bot Token

---

## Mage Commands

```bash
mage setup         # Install dev tools
mage dev           # Run bot with auto-restart
mage build         # Build binary
mage fmt / lint    # Format & lint checks
mage vulncheck     # Security check
```

---

<p align="center">
  <img src="https://github.com/dunamismax/images/blob/main/golang/discord-bots/generic-bot-gopher.png" alt="discord-bot-gopher" width="300" />
</p>

## Bot Commands

```bash
# Basic commands included
!ping                 # Check if bot is online
!help                 # Show available commands
!stats                # Display bot performance metrics

# Add your own commands by extending the command handlers
```

## Bot in Action

<p align="center">
  <img src="https://github.com/dunamismax/images/blob/main/golang/discord-bots/generic-bot-help.png" alt="Help Command Screenshot" width="500" />
  <br>
  <em>Help command showing available features</em>
</p>

<p align="center">
  <img src="https://github.com/dunamismax/images/blob/main/golang/discord-bots/generic-bot-stats.png" alt="Stats Command Screenshot" width="500" />
  <br>
  <em>Performance statistics and monitoring</em>
</p>

---

## Project Structure

The template uses a clean architecture with organized packages:

* `main.go` - Application entry point with graceful shutdown
* `discord/` - Discord client and bot logic
* `config/` - Configuration management with validation
* `metrics/` - Performance monitoring and statistics
* `logging/` - Structured logging utilities
* `errors/` - Custom error types and handling
* `magefile.go` - Build automation and development tools

Start development with `mage dev` for auto-restart functionality.

---

## Adding Your Own Commands

1. **Add command handler** in `discord/bot.go`:
```go
func (b *Bot) handleMyCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
    // Your command logic here
    embed := &discordgo.MessageEmbed{
        Title: "My Command",
        Description: "This is my custom command!",
        Color: 0x00FF00,
    }
    _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
    return err
}
```

2. **Register the command** in `registerCommands()`:
```go
b.commandHandlers["mycommand"] = b.handleMyCommand
```

3. **Update help text** to document your new command.

---

<p align="center">
  <img src="https://github.com/dunamismax/images/blob/main/golang/go-logo.png" alt="Go Logo" width="300" />
</p>

## Deployment Options

* **Single Binary** – Build with `mage build`, copy the file, and run with env vars.
* **Systemd** – Create a service to keep it running on Linux.
* **Docker** – Lightweight container build can be easily added.
* **Cloud Platforms** – Deploy to any Go-supporting platform.

---

## Configuration

Environment variables with sensible defaults:

```bash
# Required
DISCORD_TOKEN=your_bot_token_here

# Optional (with defaults)
COMMAND_PREFIX=!
LOG_LEVEL=info
DEBUG=false
BOT_NAME=discord-bot
JSON_LOGGING=false
SHUTDOWN_TIMEOUT=30s
REQUEST_TIMEOUT=30s
MAX_RETRIES=3
```

---

<p align="center">
  <a href="https://buymeacoffee.com/dunamismax" target="_blank">
    <img src="https://github.com/dunamismax/images/blob/main/golang/buy-coffee-go.gif" alt="Buy Me A Coffee" style="height: 150px !important;" />
  </a>
</p>

<p align="center">
  <a href="https://twitter.com/dunamismax" target="_blank"><img src="https://img.shields.io/badge/Twitter-%231DA1F2.svg?&style=for-the-badge&logo=twitter&logoColor=white" alt="Twitter"></a>
  <a href="https://bsky.app/profile/dunamismax.bsky.social" target="_blank"><img src="https://img.shields.io/badge/Bluesky-blue?style=for-the-badge&logo=bluesky&logoColor=white" alt="Bluesky"></a>
  <a href="https://reddit.com/user/dunamismax" target="_blank"><img src="https://img.shields.io/badge/Reddit-%23FF4500.svg?&style=for-the-badge&logo=reddit&logoColor=white" alt="Reddit"></a>
  <a href="https://discord.com/users/dunamismax" target="_blank"><img src="https://img.shields.io/badge/Discord-dunamismax-7289DA.svg?style=for-the-badge&logo=discord&logoColor=white" alt="Discord"></a>
  <a href="https://signal.me/#p/+dunamismax.66" target="_blank"><img src="https://img.shields.io/badge/Signal-dunamismax.66-3A76F0.svg?style=for-the-badge&logo=signal&logoColor=white" alt="Signal"></a>
</p>

## License

MIT – see [LICENSE](LICENSE) for details.

---

<p align="center">
  <strong>Discord Bot Template in Go</strong><br>
  <sub>DiscordGo • Mage • slog • Config • Metrics • Clean Architecture • Production Ready</sub>
</p>

---