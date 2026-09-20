## 1. Core Chat ID Retrieval Logic

- [x] 1.1 Implement `RunChatID` in `commands.go` supporting single most-recent ID, `--all`, and `--json`, verifying with unit tests
- [x] 1.2 Wire `chat-id` and `get-chat-id` subcommands into `main.go` supporting token via `--token` and positional arguments

## 2. Agent Skill & Reference Updates

- [x] 2.1 Update `.agents/skills/telegram/SKILL.md` and `.agents/skills/telegram/references/setup.md` to document `tg chat-id` and `tg get-chat-id`

## 3. Verification

- [x] 3.1 Recompile `tg` binary, test `tg chat-id --help`, and execute test suite verifying chat ID extraction
