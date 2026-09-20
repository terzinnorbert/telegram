# Telegram CLI Setup & Configuration

This guide details how to create your Telegram Bot, obtain authentication tokens, discover your chat ID, and configure the `tg` binary.

---

## 1. Create a Bot with @BotFather

1. Open Telegram and search for [`@BotFather`](https://t.me/BotFather).
2. Send `/newbot` and follow the prompts to choose:
   - A friendly name (e.g. `DevBot`)
   - A unique username ending in `bot` (e.g. `volta_dev_bot`)
3. Copy the HTTP API token provided by BotFather (looks like `123456789:ABCdefGhIJKlmNoPQRstuVWXyz`).

---

## 2. Discover Your Chat ID

Telegram bots cannot initiate a conversation with a user unprompted. You must send a message to your new bot first:

1. Open your new bot's profile in Telegram (e.g. `https://t.me/volta_dev_bot`) and press **Start** (or send `/start` or `hello`).
2. Run the `chat-id` retrieval command:
   ```bash
   # Direct retrieval with token argument or flag:
   tg chat-id --token "<your-bot-token>"
   ```
   Or set it directly into your shell environment:
   ```bash
   export TELEGRAM_CHAT_ID=$(tg chat-id --token "<your-bot-token>")
   ```
3. Alternatively, to view details of all chats:
   ```bash
   tg chats --token "<your-bot-token>"
   ```
   Output:
   ```text
   CHAT ID       TYPE     NAME / USERNAME
   987654321     private  Norty (@norty)
   ```

---

## 3. Configuration Methods

`tg` checks for configuration in the following order:
1. Command-line flags (`--token`, `--chat-id`)
2. Environment variables (`TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`)
3. Config file (`~/.config/tg/config.json`)

### Option A: Environment Variables (Recommended for CLI / CI)
Add to your `~/.bashrc`, `~/.zshrc`, or `.env`:
```bash
export TELEGRAM_BOT_TOKEN="123456789:ABCdefGhIJKlmNoPQRstuVWXyz"
export TELEGRAM_CHAT_ID="987654321"
```

### Option B: Configuration File
Create `~/.config/tg/config.json`:
```bash
mkdir -p ~/.config/tg
cat <<EOF > ~/.config/tg/config.json
{
  "bot_token": "123456789:ABCdefGhIJKlmNoPQRstuVWXyz",
  "chat_id": "987654321"
}
EOF
chmod 600 ~/.config/tg/config.json
```

---

## 4. Verifying Connectivity

Verify your bot is responding:
```bash
tg whoami
```
Output:
```text
Bot ID:       123456789
Username:     @volta_dev_bot
Name:         DevBot
```

Send a test message:
```bash
tg send "Telegram CLI successfully configured!"
```
Check your Telegram app to confirm receipt!
