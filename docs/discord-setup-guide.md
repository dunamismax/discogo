# Discord Bot Setup Guide

Complete step-by-step guide to deploying your Go Discord bot on Discord servers.

## Table of Contents

1. [Discord Developer Portal Setup](#1-discord-developer-portal-setup)
2. [Bot Configuration](#2-bot-configuration)
3. [Server Permissions Setup](#3-server-permissions-setup)
4. [Bot Deployment](#4-bot-deployment)
5. [Testing & Verification](#5-testing--verification)
6. [Troubleshooting](#6-troubleshooting)

---

## 1. Discord Developer Portal Setup

### Step 1.1: Create Discord Application

1. **Navigate to Discord Developer Portal**
   - Go to [https://discord.com/developers/applications](https://discord.com/developers/applications)
   - Log in with your Discord account

2. **Create New Application**
   - Click **"New Application"** button (top right)
   - Enter application name: `My Discord Bot` (or your preferred name)
   - Click **"Create"**

3. **Configure Application Settings**
   - **Name**: `My Discord Bot`
   - **Description**: `A Discord bot built with Go and DiscordGo`
   - **App Icon**: Upload a bot-related image (optional)

### Step 1.2: Create Bot User

1. **Navigate to Bot Section**
   - In your application, click **"Bot"** in the left sidebar

2. **Create Bot**
   - Click **"Add Bot"** button
   - Confirm by clicking **"Yes, do it!"**

3. **Configure Bot Settings**

   **Basic Information:**
   - **Username**: `MyBot` (or your preferred name)
   - **Avatar**: Upload bot avatar image (optional)

   **Authorization Flow:**
   - ✅ **Public Bot**: `ON` (allows others to invite your bot)
   - ❌ **Require OAuth2 Code Grant**: `OFF` (not needed for most bots)

   **Privileged Gateway Intents:**
   - ❌ **Presence Intent**: `OFF` (enable if your bot needs user presence data)
   - ❌ **Server Members Intent**: `OFF` (enable if your bot needs member list access)
   - ✅ **Message Content Intent**: `ON` (required to read message content)

4. **Copy Bot Token**
   - Click **"Reset Token"** to generate a new token
   - **⚠️ IMPORTANT**: Copy the token immediately and store it securely
   - **Never share this token publicly or commit it to version control**

---

## 2. Bot Configuration

### Step 2.1: Environment Setup

1. **Clone Repository**

   ```bash
   git clone https://github.com/dunamismax/discogo.git
   cd discogo
   ```

2. **Install Dependencies**

   ```bash
   go mod tidy
   go install github.com/magefile/mage@latest
   ```

3. **Create Environment File**

   ```bash
   cp .env.example .env
   ```

4. **Configure Environment Variables**

   Edit `.env` file with your settings:

   ```bash
   # Required - Your Discord bot token
   DISCORD_TOKEN=your_bot_token_here_from_step_1.2
   
   # Optional - Bot configuration
   COMMAND_PREFIX=!
   LOG_LEVEL=info
   BOT_NAME=discord-bot
   DEBUG=false
   ```

### Step 2.2: Development Setup

1. **Install Development Tools**

   ```bash
   mage setup
   ```

2. **Verify Configuration**

   ```bash
   mage status
   ```

   Expected output:

   ```
   Discord Bot Development Environment Status
   =========================================
   Go: go version go1.24+ darwin/arm64
   Environment: .env file found ✓
   Bot: discord-bot main.go found ✓
   Built binaries: None found
   ```

---

## 3. Server Permissions Setup

### Step 3.1: Generate Invite Link

1. **Navigate to OAuth2 Section**
   - In Discord Developer Portal, go to **"OAuth2"** → **"URL Generator"**

2. **Select Scopes**
   - ✅ **bot**: Required for bot functionality
   - ✅ **applications.commands**: Required for slash commands (if needed)

3. **Select Bot Permissions**

   **Essential Permissions:**
   - ✅ **Send Messages**: Post messages
   - ✅ **Embed Links**: Display rich embeds
   - ✅ **Read Message History**: Process commands

   **Optional Permissions:**
   - ✅ **Add Reactions**: React to messages
   - ✅ **Manage Messages**: Delete messages (if needed)
   - ✅ **Use External Emojis**: Enhanced formatting

4. **Copy Generated URL**
   - Copy the generated URL at the bottom
   - Example: `https://discord.com/api/oauth2/authorize?client_id=YOUR_CLIENT_ID&permissions=412384301120&scope=bot%20applications.commands`

### Step 3.2: Add Bot to Server

1. **Visit Invite URL**
   - Paste the generated URL into your browser
   - Log in to Discord if prompted

2. **Select Server**
   - Choose the Discord server where you want to add the bot
   - You must have **"Manage Server"** permission

3. **Confirm Permissions**
   - Review the permissions list
   - Click **"Authorize"**
   - Complete any CAPTCHA if presented

4. **Verify Bot Addition**
   - Check your Discord server's member list
   - The bot should appear offline (until you start it)

---

## 4. Bot Deployment

### Step 4.1: Local Development

**Start Development Server:**

```bash
mage dev
```

Expected output:

```
Starting discord-bot in development mode with auto-restart...
Press Ctrl+C to stop.
2024/08/08 12:00:00 Starting Discord Bot...
2024/08/08 12:00:00 Bot Name: discord-bot
2024/08/08 12:00:00 Command Prefix: !
2024/08/08 12:00:00 Log Level: info
2024/08/08 12:00:00 Bot is now running
```

### Step 4.2: Production Deployment

**Option A: Direct Binary Deployment**

1. **Build Production Binary**

   ```bash
   mage build
   ```

2. **Deploy Binary**

   ```bash
   # Copy binary to production server
   scp bin/discord-bot user@your-server:/opt/discord-bot/
   
   # Run with environment variables
   DISCORD_TOKEN=your_token LOG_LEVEL=info ./discord-bot
   ```

**Option B: Systemd Service (Ubuntu/Linux)**

1. **Create Service User**

   ```bash
   sudo useradd -r -s /bin/false discordbot
   sudo mkdir -p /opt/discord-bot
   sudo chown discordbot:discordbot /opt/discord-bot
   ```

2. **Deploy Binary**

   ```bash
   sudo cp bin/discord-bot /opt/discord-bot/
   sudo chmod +x /opt/discord-bot/discord-bot
   sudo chown discordbot:discordbot /opt/discord-bot/discord-bot
   ```

3. **Create Systemd Service**

   ```bash
   sudo tee /etc/systemd/system/discord-bot.service > /dev/null << EOF
   [Unit]
   Description=Discord Bot
   After=network.target
   Wants=network.target
   
   [Service]
   Type=simple
   User=discordbot
   Group=discordbot
   WorkingDirectory=/opt/discord-bot
   ExecStart=/opt/discord-bot/discord-bot
   Restart=always
   RestartSec=10
   
   # Environment variables
   Environment=DISCORD_TOKEN=your_bot_token_here
   Environment=LOG_LEVEL=info
   Environment=COMMAND_PREFIX=!
   Environment=BOT_NAME=discord-bot
   
   # Security settings
   NoNewPrivileges=true
   PrivateTmp=true
   ProtectSystem=strict
   ProtectHome=true
   ReadWritePaths=/opt/discord-bot
   
   [Install]
   WantedBy=multi-user.target
   EOF
   ```

4. **Start and Enable Service**

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable discord-bot
   sudo systemctl start discord-bot
   ```

5. **Check Service Status**

   ```bash
   sudo systemctl status discord-bot
   sudo journalctl -u discord-bot -f
   ```

---

## 5. Testing & Verification

### Step 5.1: Basic Functionality Tests

1. **Verify Bot is Online**
   - Check Discord server member list
   - Bot should show as "Online" with green status

2. **Test Basic Commands**

   ```
   # Test ping command
   !ping
   
   # Test help command
   !help
   
   # Test stats command
   !stats
   ```

3. **Expected Response Format**
   - Rich embed responses
   - Proper error handling
   - Quick response times

### Step 5.2: Error Handling Tests

1. **Test Invalid Command**

   ```
   !invalidcommand
   ```

   Expected: Error message with help suggestion

2. **Test Empty Command**

   ```
   !
   ```

   Expected: No response (ignored)

---

## 6. Troubleshooting

### Common Issues

#### Bot Appears Offline

**Symptoms**: Bot shows as offline in Discord

**Solutions**:

1. **Check Token**: Verify `DISCORD_TOKEN` is correct in `.env`
2. **Check Intents**: Ensure "Message Content Intent" is enabled
3. **Check Logs**: Look for authentication errors

   ```bash
   mage dev
   # or
   sudo journalctl -u discord-bot -f
   ```

#### Bot Doesn't Respond to Commands

**Symptoms**: Bot is online but ignores commands

**Possible Causes & Solutions**:

1. **Missing Message Content Intent**
   - Go to Discord Developer Portal → Bot → Privileged Gateway Intents
   - Enable "Message Content Intent"

2. **Wrong Command Prefix**
   - Check `COMMAND_PREFIX` in `.env` (default: `!`)
   - Test with correct prefix: `!ping`

3. **Insufficient Permissions**
   - Check bot has "Send Messages" and "Embed Links" permissions
   - Re-invite bot with correct permissions

### Log Analysis

**Viewing Logs**:

```bash
# Development mode
mage dev

# Production systemd
sudo journalctl -u discord-bot -f

# Docker
docker logs -f container_name
```

**Important Log Messages**:

- ✅ `Bot is now running` - Successful startup
- ❌ `invalid authentication token` - Wrong Discord token
- ❌ `missing access` - Insufficient permissions

### Getting Help

1. **Check Repository Issues**: [GitHub Issues](https://github.com/dunamismax/discogo/issues)
2. **Discord Developer Documentation**: [Discord API Docs](https://discord.com/developers/docs)

---

## Quick Reference Commands

```bash
# Development
mage dev                     # Start with auto-restart
mage status                  # Check environment
mage build                   # Build production binary

# Production Management
sudo systemctl status discord-bot
sudo systemctl restart discord-bot
sudo journalctl -u discord-bot -f

# Testing Commands
!ping                        # Test basic functionality
!help                        # Show available commands
!stats                       # Show bot statistics
```

---

**Success!** Your Discord bot should now be running and responding to commands in your Discord server.