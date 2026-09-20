## ADDED Requirements

### Requirement: Direct chat ID retrieval
The CLI SHALL provide a `chat-id` command (with alias `get-chat-id`) that accepts a bot token as a parameter and outputs the most recent active chat ID.

#### Scenario: Retrieve single chat ID with token flag
- **WHEN** user executes `tg chat-id --token "<bot-token>"` and updates exist
- **THEN** the CLI SHALL print the single most recent active chat ID directly to stdout and exit with status 0

#### Scenario: Retrieve single chat ID with positional argument
- **WHEN** user executes `tg chat-id "<bot-token>"` and updates exist
- **THEN** the CLI SHALL resolve the token from the positional argument and print the most recent active chat ID to stdout

#### Scenario: Retrieve all chat IDs
- **WHEN** user executes `tg chat-id --token "<bot-token>" --all`
- **THEN** the CLI SHALL print all discovered unique chat IDs, one per line, to stdout

#### Scenario: Retrieve chat IDs as JSON
- **WHEN** user executes `tg chat-id --token "<bot-token>" --json`
- **THEN** the CLI SHALL output valid JSON containing discovered chat details (`id`, `type`, `name`, `username`)

#### Scenario: No chats found
- **WHEN** user executes `tg chat-id --token "<bot-token>"` but no incoming updates exist for the bot
- **THEN** the CLI SHALL output an actionable error message to stderr and exit with a non-zero status code
