## 1. Toolchain & Project Scaffolding

- [x] 1.1 Set up Go environment, initialize `go.mod`, and verify compilation with `go build`
- [x] 1.2 Implement CLI flag parsing and configuration loader (flags > env vars > config file) and verify precedence with unit tests

## 2. Telegram API Client & Discovery

- [x] 2.1 Implement HTTP client wrapper for Telegram Bot API with timeout and retry handling for `getMe`, `sendMessage`, `sendDocument`, and `getUpdates`
- [x] 2.2 Implement `tg whoami` to verify bot credentials and `tg chats` to list recent chats and IDs from incoming updates

## 3. Message & Attachment Sending

- [x] 3.1 Implement `tg send [text]` supporting command arguments, standard input piping, and parse modes (`markdown`, `html`, `plain`), verifying output with automated tests
- [x] 3.2 Implement `tg send -f <path> -c <caption>` for uploading documents and files via multipart/form-data, verifying file existence and size checks

## 4. Stateless Message Fetching

- [x] 4.1 Implement `tg fetch` with `--limit` and `--chat-id` filters for human-readable terminal output
- [x] 4.2 Implement `tg fetch --json` structured output matching the specification for automated agent consumption

## 5. Agent Skill & Reference Guides

- [x] 5.1 Create `.agents/skills/telegram/SKILL.md` containing agent workflow rules, quick commands, and routing
- [x] 5.2 Create `.agents/skills/telegram/references/sending.md` with guidelines for file attachments, formatting, and stdin piping
- [x] 5.3 Create `.agents/skills/telegram/references/fetching.md` detailing stateless polling and JSON schema parsing
- [x] 5.4 Create `.agents/skills/telegram/references/setup.md` covering BotFather setup, chat ID discovery, and environment variables

## 6. End-to-End Verification

- [x] 6.1 Compile the standalone `tg` binary, test all `--help` flags, and verify binary size and execution
- [x] 6.2 Execute automated test suite verifying JSON serialization, multipart uploads, and error handling
