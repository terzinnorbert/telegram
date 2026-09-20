---
name: telegram
description: Send and fetch Telegram messages and upload file attachments (such as Markdown reports, logs, and code) via the `tg` CLI tool.
allowed-tools: Bash(tg:*)
metadata:
  author: volta
  version: "1.0.0"
---

# Telegram CLI Agent Skill

Use this skill whenever you need to interact with Telegram: sending status notifications, delivering generated files or markdown reports to a user or channel, or reading incoming instructions and feedback.

The underlying CLI tool is `tg` (a single standalone Go binary).

## Quick Decision Tree

- **Want to send a brief status or notification?**
  ```bash
  tg send "Build completed successfully!"
  ```

- **Want to send a generated Markdown file or document?**
  ```bash
  tg send -f path/to/document.md -c "Sprint Report"
  ```

- **Want to pipe output from another command into Telegram?**
  ```bash
  rtk git diff | tg send
  ```

- **Want to read recent incoming messages?**
  ```bash
  tg fetch --limit 5 --json
  ```

- **Want to discover chat IDs or configure credentials?**
  ```bash
  # Directly retrieve active chat ID (great for scripts and env assignment)
  tg chat-id --token "<token>"

  # List all discovered chats in a table
  tg chats
  tg whoami
  ```

---

## Detailed References

For in-depth guides and advanced usage, read the following reference docs:

- [Sending Guide](references/sending.md): Markdown formatting, special character escaping, file attachment limits, captions, and stdin streaming.
- [Fetching Guide](references/fetching.md): Stateless polling, JSON output schema, parsing incoming text, and filtering by chat ID.
- [Setup & Configuration](references/setup.md): Obtaining a BotFather token, configuring environment variables (`TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`), and testing connectivity.

---

## Safety & Operational Rules

1. **Verify Credentials First**: If `TELEGRAM_BOT_TOKEN` or `TELEGRAM_CHAT_ID` are not set, consult [Setup & Configuration](references/setup.md).
2. **Prefer Document Upload for Long Content**: If sending reports, logs, or code longer than ~20 lines, write them to a markdown file and send via `tg send -f <path> -c "..."` rather than sending a huge text wall.
3. **Parse Mode Caveats**: By default, `tg send` uses Telegram Markdown mode. If your text contains raw unescaped characters (like random underscores or brackets in log dumps), use `--parse-mode plain` or upload as a document with `-f`.
