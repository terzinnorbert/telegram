## Why

When configuring bots or writing scripts and agent tool workflows, extracting a chat ID from recent updates currently requires running `tg chats` and manually parsing a tabular output. Users and automated agents need a direct, script-friendly command (`tg chat-id`) that accepts a bot token as a parameter and directly outputs the target chat ID.

## What Changes

- Add a dedicated `chat-id` subcommand (with alias `get-chat-id`) to the `tg` binary.
- Support passing the token via `--token` flag, positional argument `tg chat-id <token>`, or the `TELEGRAM_BOT_TOKEN` environment variable.
- In default mode, output the single most recent active chat ID directly to stdout for shell assignment (e.g., `export TELEGRAM_CHAT_ID=$(tg chat-id --token "$TOKEN")`).
- Support `--all` to print all discovered unique chat IDs, and `--json` for machine-readable JSON output.
- Return a descriptive error to stderr if no messages have been sent to the bot yet.
- Update agent skill documentation to guide agents on using `tg chat-id --token <token>` for automated chat discovery.

## Capabilities

### Modified Capabilities
- `telegram-cli`: Add direct chat ID retrieval requirement and command options.
- `telegram-agent-skill`: Document the `chat-id` command in the agent skill setup reference.

## Impact

- CLI: Adds `tg chat-id` and `tg get-chat-id` commands.
- Skill: Updates setup instructions in `.agents/skills/telegram/`.
- Zero breaking changes to existing commands (`tg chats`, `tg send`, `tg fetch`).
