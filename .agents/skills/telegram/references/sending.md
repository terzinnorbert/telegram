# Sending Messages & Attachments Guide

This guide details how agents and scripts can transmit text, format markdown, pipe command output, and upload documents using `tg send`.

---

## 1. Sending Text Messages

### Direct Arguments
Pass the message text directly to `tg send`:
```bash
tg send "Deployment finished with exit status 0"
```

### Specifying Target Chat ID
If `TELEGRAM_CHAT_ID` is set in the environment or config file, `--chat-id` is optional. To target a specific user or group:
```bash
tg send --chat-id 12345678 "Message for specific user"
```

### Parse Modes & Formatting
Telegram supports formatted text. Specify `--parse-mode` as `markdown` (default), `html`, or `plain`.

- **Markdown mode** (default):
  ```bash
  tg send "*Bold Title*\n- Item 1\n- Item 2\n\`inline code\`"
  ```
- **HTML mode**:
  ```bash
  tg send --parse-mode html "<b>Important:</b> <code>System online</code>"
  ```
- **Plain text mode**:
  Use `plain` mode when sending raw stack traces, log excerpts, or strings with arbitrary special characters (`_`, `*`, `[`, `]`) that could break Telegram Markdown parsing:
  ```bash
  tg send --parse-mode plain "Error in function_name_with_underscores[0]"
  ```

---

## 2. Standard Input Piping

You can pipe data directly from command outputs into `tg send`:
```bash
# Pipe git status
rtk git status | tg send

# Pipe logs
journalctl -u my-service -n 20 --no-pager | tg send --parse-mode plain
```

---

## 3. Uploading Files & Document Attachments

When the user asks you to "send the generated markdown to telegram", use the `-f` (or `--file`) flag:

```bash
tg send -f report.md -c "Sprint Analysis Report"
```

### Options:
- `-f <path>`: Local file path to upload.
- `-c <text>`: Optional caption accompanying the document.
- `--chat-id <id>`: Target recipient or channel.
- `--parse-mode <mode>`: Formatting mode for the caption (default: `markdown`).

### Special Stdin Upload:
To upload piped content as a file document attachment instead of an inline text message, use `-f -`:
```bash
rtk git diff | tg send -f - -c "Current unstaged git diff"
```

### Attachment Constraints:
- Maximum file size: **50 MB** (Telegram Bot API limit).
- The binary automatically inspects file size before uploading and returns an error if exceeded.
