# Fetching Messages Guide

This guide details how agents and scripts can retrieve incoming messages statelessly from Telegram using `tg fetch`.

---

## 1. Stateless Polling Model

`tg fetch` connects to Telegram's `getUpdates` API endpoint. It operates **statelessly** by default, allowing agents to query context on-demand without managing local cursor state files.

### Basic Human-Readable Retrieval:
```bash
tg fetch --limit 5
```

Sample output:
```text
[2026-09-20T15:00:00Z] @norty (Chat 12345678 | Msg 42):
Please deploy the new migration

[2026-09-20T15:05:00Z] @norty (Chat 12345678 | Msg 43):
Also attach the coverage report
```

---

## 2. Structured JSON Output for Agents

When calling `tg fetch` from an agent or automated script, always use `--json`:
```bash
tg fetch --limit 5 --json
```

### JSON Schema:
The output is a JSON array of message objects:
```json
[
  {
    "update_id": 987654,
    "message_id": 42,
    "chat_id": 12345678,
    "chat_title": "",
    "sender": "@norty",
    "date": "2026-09-20T15:00:00Z",
    "text": "Please deploy the new migration"
  }
]
```

### Fields:
- `update_id` (number): Telegram internal update identifier.
- `message_id` (number): Message identifier within the chat.
- `chat_id` (number): Numeric identifier of the chat (can be positive for private chats, negative for groups/channels).
- `chat_title` (string): Title of group/channel (if applicable).
- `sender` (string): `@username` or first/last name of the sender.
- `date` (string): ISO 8601 UTC timestamp.
- `text` (string): Text content or caption of the message.

---

## 3. Filtering Messages by Chat ID

To isolate messages coming from a specific user or group chat:
```bash
tg fetch --chat-id 12345678 --limit 10 --json
```
Only messages whose `chat_id` equals `12345678` will be returned.

---

## 4. Manual Offset Progression (Optional)

If an agent needs to retrieve updates newer than a specific `update_id`:
```bash
tg fetch --offset 987655 --limit 10 --json
```
This queries updates with `update_id >= 987655`.
