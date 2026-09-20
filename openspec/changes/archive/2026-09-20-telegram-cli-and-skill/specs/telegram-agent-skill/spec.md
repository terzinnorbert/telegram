## Purpose

Defines an agent skill and specialized reference documentation enabling autonomous AI agents to interact with Telegram via the `tg` CLI tool.

## ADDED Requirements

### Requirement: Root skill definition
The skill SHALL provide a root `SKILL.md` in `.agents/skills/telegram/` specifying tool capabilities, command invocations, triage rules, and pointers to detailed subpages.

#### Scenario: Agent checks available Telegram operations
- **WHEN** an agent reads `.agents/skills/telegram/SKILL.md`
- **THEN** it SHALL find clear instructions and examples for sending text, uploading files, fetching incoming messages, and locating subpages

### Requirement: Sending reference subpage
The skill SHALL include a reference guide at `.agents/skills/telegram/references/sending.md` detailing sending patterns, document uploads, captions, markdown formatting rules, and standard input piping.

#### Scenario: Agent needs to send a generated markdown file
- **WHEN** an agent inspects `.agents/skills/telegram/references/sending.md`
- **THEN** it SHALL find syntax guidelines for `tg send -f <path> -c "<caption>"` and file size/MIME handling

### Requirement: Fetching reference subpage
The skill SHALL include a reference guide at `.agents/skills/telegram/references/fetching.md` describing stateless polling, JSON output parsing, and chat filtering.

#### Scenario: Agent queries incoming instructions
- **WHEN** an agent inspects `.agents/skills/telegram/references/fetching.md`
- **THEN** it SHALL find instructions on invoking `tg fetch --limit N --json` and the JSON schema specification for processing the response

### Requirement: Setup reference subpage
The skill SHALL include a reference guide at `.agents/skills/telegram/references/setup.md` guiding users and agents through BotFather bot registration, chat ID discovery via `tg chats`, and environment configuration.

#### Scenario: Agent assists user with initial bot configuration
- **WHEN** an agent inspects `.agents/skills/telegram/references/setup.md`
- **THEN** it SHALL find step-by-step instructions for creating a bot with @BotFather and setting `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID`
