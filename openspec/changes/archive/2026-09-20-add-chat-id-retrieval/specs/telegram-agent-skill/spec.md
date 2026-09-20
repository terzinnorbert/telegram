## MODIFIED Requirements

### Requirement: Setup reference subpage
The skill SHALL include a reference guide at `.agents/skills/telegram/references/setup.md` guiding users and agents through BotFather bot registration, direct chat ID discovery via `tg chat-id` (and `tg chats`), and environment configuration.

#### Scenario: Agent assists user with initial bot configuration
- **WHEN** an agent inspects `.agents/skills/telegram/references/setup.md`
- **THEN** it SHALL find step-by-step instructions for creating a bot with @BotFather, retrieving chat ID directly with `tg chat-id --token <token>`, and setting `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID`
