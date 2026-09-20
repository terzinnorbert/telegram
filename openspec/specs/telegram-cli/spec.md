# telegram-cli Specification

## Purpose

Provides a fast, dependency-free command-line interface in Go for sending messages, uploading documents and attachments, and querying incoming updates via the Telegram Bot API.

## Requirements

### Requirement: Send text messages
The CLI SHALL allow sending text messages to a Telegram chat, accepting message text either as a command argument or via standard input.

#### Scenario: Send message as argument
- **WHEN** user executes `tg send "Deploy completed"`
- **THEN** the message SHALL be transmitted to the configured default Telegram chat ID and exit with status 0 upon success

#### Scenario: Send message via stdin
- **WHEN** user executes `echo "Build logs" | tg send`
- **THEN** the content read from stdin SHALL be sent to the configured Telegram chat

#### Scenario: Send message to custom chat ID
- **WHEN** user executes `tg send --chat-id 987654321 "Direct notification"`
- **THEN** the message SHALL be sent specifically to chat ID `987654321` overriding the default chat ID

### Requirement: Formatting parse modes
The CLI SHALL support `markdown`, `html`, and `plain` text parse modes when sending messages.

#### Scenario: Default Markdown formatting
- **WHEN** user executes `tg send "*Bold* and _Italic_"` without specifying `--parse-mode`
- **THEN** the CLI SHALL format and transmit the message using Telegram's Markdown parse mode

#### Scenario: HTML parse mode
- **WHEN** user executes `tg send --parse-mode html "<b>Bold text</b>"`
- **THEN** the CLI SHALL send the message using Telegram's HTML parse mode

### Requirement: File and document attachments
The CLI SHALL support uploading files and documents (such as Markdown files, text logs, images, and archives) with an optional caption.

#### Scenario: Send file attachment with caption
- **WHEN** user executes `tg send -f report.md -c "Release summary"`
- **THEN** the file `report.md` SHALL be uploaded via Telegram's `sendDocument` endpoint with the caption `"Release summary"`

#### Scenario: Missing file error
- **WHEN** user executes `tg send -f non_existent_file.md`
- **THEN** the CLI SHALL terminate with a non-zero exit code and output a descriptive error message to stderr

### Requirement: Stateless message fetching
The CLI SHALL support fetching recent incoming messages via the Telegram Bot API `getUpdates` endpoint in a stateless manner without requiring a local offset cursor.

#### Scenario: Fetch recent messages in human-readable format
- **WHEN** user executes `tg fetch --limit 5`
- **THEN** up to 5 of the most recent incoming messages SHALL be printed to stdout with sender, date, and message text

#### Scenario: Fetch recent messages as JSON
- **WHEN** user executes `tg fetch --limit 5 --json`
- **THEN** the output SHALL be valid JSON containing an array of message objects with `update_id`, `message_id`, `chat_id`, `sender`, `date`, and `text` fields

#### Scenario: Filter messages by chat ID
- **WHEN** user executes `tg fetch --chat-id 12345678 --json`
- **THEN** only messages matching chat ID `12345678` SHALL be returned in the output array

### Requirement: Bot identity and chat discovery
The CLI SHALL provide helper commands to verify bot credentials and discover available chat IDs from recent interactions.

#### Scenario: Verify bot identity
- **WHEN** user executes `tg whoami`
- **THEN** the CLI SHALL query the `getMe` endpoint and output the bot's ID, username, and display name

#### Scenario: List active chats
- **WHEN** user executes `tg chats`
- **THEN** the CLI SHALL scan recent updates and display a deduplicated table of chat IDs, usernames, and titles

### Requirement: Configuration and credentials
The CLI SHALL read configuration credentials with precedence: command-line flags > environment variables (`TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`) > configuration file (`~/.config/tg/config.json`).

#### Scenario: Missing bot token error
- **WHEN** user runs `tg send "test"` without token in flags, environment variables, or config file
- **THEN** the CLI SHALL exit with a non-zero status code and inform the user that `TELEGRAM_BOT_TOKEN` is required

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

