# Claude Code Guidelines for `tg`

This repository provides `tg`, a standalone Go CLI tool for sending Telegram messages, uploading file attachments (e.g. Markdown reports, logs), and statelessly querying incoming messages.

## Build & Test Commands

- **Build binary:** `make build` (compiles stripped binary to `bin/tg`)
- **Run tests:** `make test`
- **Install globally:** `make install` (installs to `$GOPATH/bin/tg`)
- **Clean build artifacts:** `make clean`

## Telegram CLI Tool Usage

Whenever the user asks you to send notifications, upload reports, or inspect Telegram messages, use `tg` (or `./bin/tg`):

```bash
# 1. Send text message
tg send "Task completed successfully"

# 2. Upload file/document (e.g. markdown summary, report, log)
tg send -f <path> -c "<caption>"

# 3. Retrieve chat ID directly
tg chat-id --token "<bot-token>"

# 4. Fetch incoming messages as structured JSON
tg fetch --limit 5 --json

# 5. Pipe command output
git status | tg send
```

## Detailed Skills & References

Full agent skill specifications and deep references are available at:
- **Skill Overview:** [`.agents/skills/telegram/SKILL.md`](.agents/skills/telegram/SKILL.md)
- **Sending & Attachments:** [`.agents/skills/telegram/references/sending.md`](.agents/skills/telegram/references/sending.md)
- **Message Fetching:** [`.agents/skills/telegram/references/fetching.md`](.agents/skills/telegram/references/fetching.md)
- **Setup & Credentials:** [`.agents/skills/telegram/references/setup.md`](.agents/skills/telegram/references/setup.md)
