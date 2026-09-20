## Why

There is currently no unified, lightweight tool for programmatic messaging with Telegram that serves both terminal CLI workflows and autonomous AI agents. Developers and AI agents need a fast, dependency-free binary to send text, pipe command outputs, attach markdown/data files, and query incoming messages without running a daemon or managing heavy MTProto sessions.

## What Changes

- Create a small, standalone Go binary (`tg`) using the standard Telegram Bot API.
- Support sending plain text, Markdown, and HTML formatted messages with automatic parse-mode detection and stdin streaming.
- Support uploading files and document attachments (markdown documents, logs, images) with optional captions.
- Support stateless fetching of incoming updates (`tg fetch`) with human-readable and structured `--json` output for automated consumption.
- Provide a helper command (`tg chats`) to discover chat IDs and sender details from recent incoming updates to simplify initial configuration.
- Support configuration via environment variables (`TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`) and fallback configuration files.
- Create an agent skill (`.agents/skills/telegram/SKILL.md`) along with specialized subpage references (`references/sending.md`, `references/fetching.md`, and `references/setup.md`) enabling AI coding assistants to autonomously fetch and dispatch messages and files.

## Capabilities

### New Capabilities
- `telegram-cli`: Command-line tool implemented in Go for sending messages, uploading documents/attachments, and querying updates via the Telegram Bot API.
- `telegram-agent-skill`: Agent skill specification and reference guides enabling AI agents to interact with Telegram via the `tg` binary.

### Modified Capabilities
<!-- No existing capabilities modified -->

## Impact

- **CLI / Tools**: Introduces `tg` executable into the repository build workflow and PATH.
- **Skills**: Introduces `.agents/skills/telegram/` skill definition with subpages for agent tool usage.
- **Dependencies**: Pure Go standard library (`net/http`, `mime/multipart`, `encoding/json`, `flag`), requiring zero external dependencies.
