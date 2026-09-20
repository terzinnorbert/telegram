## Context

See proposal.md. The current CLI provides `tg chats` which renders an ASCII table of all discovered chats. While useful for human inspection, scripts, CI environments, and autonomous AI agents require a clean, single-value command (`tg chat-id`) to retrieve the active chat ID directly.

## Goals / Non-Goals

**Goals:**
- Provide `tg chat-id` (and alias `tg get-chat-id`).
- Accept the bot token as a parameter via `--token`, positional argument, or `TELEGRAM_BOT_TOKEN` environment variable.
- By default, print only the raw numeric ID of the most recent chat directly to stdout, exiting 0.
- Support `--all` (list all unique IDs) and `--json` (structured list).
- Exit non-zero with actionable stderr instructions if no chats have contacted the bot.

**Non-Goals:**
- Interactive prompts or interactive TUI selection.

## Decisions

### Decision: Raw Single-Value Output as Default
- **Choice**: Default output of `tg chat-id` is only the raw integer chat ID (e.g., `123456789\n`) without table headers or decorative prefixes.
- **Rationale**: Enables direct shell assignment (`CHAT_ID=$(tg chat-id --token "$TOKEN")`) and painless parsing for AI agents invoking bash tools.
- **Alternatives Considered**: Always returning tabular or JSON output (requires extra parsing steps like `jq` or `awk`).

### Decision: Most Recent Chat Selection
- **Choice**: When multiple chats exist, default mode selects the chat associated with the latest update (highest `date` or `update_id`).
- **Rationale**: The user setting up the bot typically just sent a message to the bot right before running the command.

### Decision: Token Parameter Ergonomics
- **Choice**: Allow token to be provided either via `--token <val>` or as a positional argument: `tg chat-id <token>`.
- **Rationale**: Matches user preference for quick one-liner execution: `tg chat-id 123456:ABC...`.

## Risks / Trade-offs

- **[Risk] Bot has received zero messages**: New bots have empty `getUpdates` until a user messages them.
  - *Mitigation*: Output a clear instruction to stderr: `"No recent chats found. Open Telegram, send a message (like /start) to your bot, then rerun this command."` and exit with status 1.
