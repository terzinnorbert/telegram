# `tg` - Lightweight Telegram CLI & AI Agent Tool

[![Go Version](https://img.shields.io/badge/go-1.19%2B-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A fast, single-binary Telegram CLI tool designed for developers, CI/CD pipelines, and autonomous AI coding agents. 

Send messages, stream piped output, upload file attachments (like Markdown reports or build logs), and statelessly query incoming messages with zero external dependencies.

---

## Highlights

- **Zero External Dependencies**: Built entirely with Go standard library (`net/http`, `mime/multipart`, `encoding/json`).
- **Tiny & Standalone**: Produces a compact static binary (~4.7MB stripped) with sub-15ms startup time.
- **File Attachments**: Upload Markdown documents, logs, images, and archives with optional captions.
- **Direct Chat ID Retrieval**: Automatically extract your active Telegram chat ID with `tg chat-id <token>`.
- **Stateless Message Fetching**: Poll recent messages on-demand with human-friendly tables or structured `--json` output.
- **Agent Skill Included**: Comes with ready-to-use Agent Skills (`.agents/skills/telegram/`) for AI pair programmers.

---

## Installation

### From Source (using Make)

```bash
git clone https://github.com/terzinnorbert/telegram.git
cd telegram
make build
```
The compiled binary will be placed at `bin/tg`.

### Install to System

```bash
make install
# Installs to $GOPATH/bin/tg (usually ~/go/bin/tg)
```

Or copy manually:
```bash
sudo cp bin/tg /usr/local/bin/
```

---

## Quickstart (2-Minute Setup)

### 1. Create a Bot
1. Open Telegram and search for [`@BotFather`](https://t.me/BotFather).
2. Send `/newbot` and follow the prompts to choose a name and username (e.g. `my_ci_bot`).
3. Copy the HTTP API token provided by BotFather.

### 2. Retrieve Your Chat ID
Telegram bots cannot message a user first. You must send a message to your new bot:
1. Open your bot in Telegram and press **Start** (or send a message like `/start` or `hello`).
2. Run the `chat-id` retrieval command:
   ```bash
   ./bin/tg chat-id "<your-bot-token>"
   ```
   *Output:*
   ```text
   987654321
   ```

### 3. Configure Credentials
Set your credentials in your shell or `.env`:
```bash
export TELEGRAM_BOT_TOKEN="<your-bot-token>"
export TELEGRAM_CHAT_ID="<your-chat-id>"
```

*(Alternatively, create a config file at `~/.config/tg/config.json`:)*
```json
{
  "bot_token": "123456789:ABCdefGhIJKlmNoPQRstuVWXyz",
  "chat_id": "987654321"
}
```

Verify your setup:
```bash
tg whoami
tg send "Hello from tg CLI!"
```

---

## CLI Commands & Usage

### 1. Sending Messages & Text

```bash
# Direct text message
tg send "Deployment completed with exit code 0"

# Target a specific chat ID (overrides default)
tg send --chat-id 987654321 "Direct notification"

# Pipe output from terminal commands
git status | tg send
cat summary.md | tg send

# Formatting modes (markdown, html, plain)
tg send "*Bold Title*\n- Item 1\n- Item 2"
tg send --parse-mode html "<b>Important:</b> <code>Service active</code>"
tg send --parse-mode plain "Error in file_with_raw_underscores[0]"
```

### 2. Uploading Files & Document Attachments

```bash
# Upload a markdown report with a caption
tg send -f report.md -c "Release Notes & Sprint Summary"

# Upload log files or images
tg send -f /var/log/build.log -c "Failed build logs"

# Upload piped content as a file document
git diff | tg send -f - -c "Git diff attachment"
```

### 3. Stateless Message Fetching

Retrieve recent incoming updates without managing local state files or cursors:

```bash
# Human-readable output
tg fetch --limit 5

# Structured JSON output (ideal for scripts & AI tools)
tg fetch --limit 5 --json

# Filter messages by specific chat ID
tg fetch --chat-id 987654321 --json
```

**Sample JSON Output:**
```json
[
  {
    "update_id": 987654,
    "message_id": 42,
    "chat_id": 987654321,
    "chat_title": "",
    "sender": "@norty",
    "date": "2026-09-20T15:00:00Z",
    "text": "Please run the migration"
  }
]
```

### 4. Chat Discovery & Bot Info

```bash
# Retrieve most recent chat ID (great for shell assignment)
export TELEGRAM_CHAT_ID=$(tg chat-id --token "$BOT_TOKEN")

# Display a table of all active chats
tg chats

# Check bot identity and verify connectivity
tg whoami
```

---

## Configuration Precedence

Credentials and settings are resolved in the following priority order:
1. **Command-line flags** (`--token`, `--chat-id`, `--config`)
2. **Environment variables** (`TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`, `TG_CONFIG`)
3. **Configuration file** (`~/.config/tg/config.json`)

---

## AI Agent Skill Integration

This repository includes a native agent skill located in [`.agents/skills/telegram/`](.agents/skills/telegram/).

When pairing with AI coding assistants (like Antigravity, Claude Code, or Cursor):
- You can instruct the agent: *"Send the generated markdown summary to Telegram"*.
- The agent reads [`.agents/skills/telegram/SKILL.md`](.agents/skills/telegram/SKILL.md) and delegates to `tg send -f ...`.
- Deep reference pages are provided in `.agents/skills/telegram/references/`:
  - [`sending.md`](.agents/skills/telegram/references/sending.md): Formatting, escaping, document limits, and stdin streaming.
  - [`fetching.md`](.agents/skills/telegram/references/fetching.md): Polling updates and JSON schemas.
  - [`setup.md`](.agents/skills/telegram/references/setup.md): Step-by-step BotFather guide.

---

## Development & Testing

```bash
# Run test suite
make test

# Build release binary
make build

# Clean artifacts
make clean
```

---

## License

[MIT](LICENSE) © 2026
