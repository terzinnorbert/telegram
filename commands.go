package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

// RunWhoami executes the `whoami` command.
func RunWhoami(ctx context.Context, client *TelegramClient, jsonOutput bool) error {
	user, err := client.GetMe(ctx)
	if err != nil {
		return err
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(user)
	}

	fmt.Printf("Bot ID:       %d\n", user.ID)
	if user.Username != "" {
		fmt.Printf("Username:     @%s\n", user.Username)
	}
	fmt.Printf("Name:         %s\n", user.FirstName)
	return nil
}

// RunChats extracts and lists unique chats discovered in recent updates.
func RunChats(ctx context.Context, client *TelegramClient, jsonOutput bool) error {
	updates, err := client.GetUpdates(ctx, 0, 100)
	if err != nil {
		return err
	}

	chatsMap := make(map[int64]ChatSummary)
	for _, u := range updates {
		msg := u.EffectiveMessage()
		if msg == nil {
			continue
		}
		c := msg.Chat
		title := c.Title
		if title == "" {
			title = c.FirstName
			if c.LastName != "" {
				title += " " + c.LastName
			}
		}
		chatsMap[c.ID] = ChatSummary{
			ID:       c.ID,
			Type:     c.Type,
			Title:    title,
			Username: c.Username,
		}
	}

	var chatList []ChatSummary
	for _, c := range chatsMap {
		chatList = append(chatList, c)
	}
	sort.Slice(chatList, func(i, j int) bool {
		return chatList[i].ID < chatList[j].ID
	})

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(chatList)
	}

	if len(chatList) == 0 {
		fmt.Println("No recent chats found. Send a message to your bot on Telegram first, then rerun 'tg chats'.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "CHAT ID\tTYPE\tNAME / USERNAME")
	for _, c := range chatList {
		ident := c.Title
		if c.Username != "" {
			if ident != "" {
				ident += " (@" + c.Username + ")"
			} else {
				ident = "@" + c.Username
			}
		}
		fmt.Fprintf(w, "%d\t%s\t%s\n", c.ID, c.Type, ident)
	}
	return w.Flush()
}

// RunSend dispatches text or document attachments.
func RunSend(ctx context.Context, client *TelegramClient, chatID string, filePath string, caption string, parseMode string, args []string, stdin io.Reader) error {
	if chatID == "" {
		return fmt.Errorf("chat ID is required. Set TELEGRAM_CHAT_ID or pass --chat-id")
	}

	// Document attachment upload
	if filePath != "" {
		var fileReader io.Reader
		var filename string

		if filePath == "-" {
			fileReader = stdin
			filename = "attachment.txt"
		} else {
			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open attachment file %s: %w", filePath, err)
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				return fmt.Errorf("failed to stat attachment file: %w", err)
			}
			if fi.Size() > 50*1024*1024 {
				return fmt.Errorf("file size %d bytes exceeds Telegram Bot API 50MB limit", fi.Size())
			}

			fileReader = f
			filename = fi.Name()
		}

		msg, err := client.SendDocument(ctx, chatID, filename, fileReader, caption, parseMode)
		if err != nil {
			return err
		}
		fmt.Printf("File sent successfully (Message ID: %d)\n", msg.MessageID)
		return nil
	}

	// Text message
	text := ""
	if len(args) > 0 {
		for i, a := range args {
			if i > 0 {
				text += " "
			}
			text += a
		}
	} else {
		// Read from stdin
		data, err := io.ReadAll(stdin)
		if err != nil {
			return fmt.Errorf("failed to read text from stdin: %w", err)
		}
		text = string(data)
	}

	if text == "" {
		return fmt.Errorf("message text is empty. Provide text argument or pipe via stdin")
	}

	msg, err := client.SendMessage(ctx, chatID, text, parseMode)
	if err != nil {
		return err
	}
	fmt.Printf("Message sent successfully (Message ID: %d)\n", msg.MessageID)
	return nil
}

// RunFetch performs stateless retrieval of updates.
func RunFetch(ctx context.Context, client *TelegramClient, limit int, offset int64, chatIDFilter int64, jsonOutput bool) error {
	updates, err := client.GetUpdates(ctx, offset, limit)
	if err != nil {
		return err
	}

	var jsonMessages []*JSONMessage
	for _, u := range updates {
		jm := u.ToJSONMessage()
		if jm == nil {
			continue
		}
		if chatIDFilter != 0 && jm.ChatID != chatIDFilter {
			continue
		}
		jsonMessages = append(jsonMessages, jm)
	}

	if jsonOutput {
		if jsonMessages == nil {
			jsonMessages = []*JSONMessage{}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(jsonMessages)
	}

	if len(jsonMessages) == 0 {
		fmt.Println("No messages found.")
		return nil
	}

	for _, m := range jsonMessages {
		fmt.Printf("[%s] %s (Chat %d | Msg %d):\n%s\n\n", m.Date, m.Sender, m.ChatID, m.MessageID, m.Text)
	}
	return nil
}

// RunChatID directly outputs the active chat ID(s).
func RunChatID(ctx context.Context, client *TelegramClient, all bool, jsonOutput bool) error {
	updates, err := client.GetUpdates(ctx, 0, 100)
	if err != nil {
		return err
	}

	chatsMap := make(map[int64]ChatSummary)
	var latestChatID int64
	var found bool

	for _, u := range updates {
		msg := u.EffectiveMessage()
		if msg == nil {
			continue
		}
		c := msg.Chat
		title := c.Title
		if title == "" {
			title = c.FirstName
			if c.LastName != "" {
				title += " " + c.LastName
			}
		}
		chatsMap[c.ID] = ChatSummary{
			ID:       c.ID,
			Type:     c.Type,
			Title:    title,
			Username: c.Username,
		}
		latestChatID = c.ID
		found = true
	}

	if !found {
		return fmt.Errorf("no recent chats found. Open Telegram, send a message (like /start) to your bot, then rerun this command")
	}

	var chatList []ChatSummary
	for _, c := range chatsMap {
		chatList = append(chatList, c)
	}
	sort.Slice(chatList, func(i, j int) bool {
		return chatList[i].ID < chatList[j].ID
	})

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(chatList)
	}

	if all {
		for _, c := range chatList {
			fmt.Printf("%d\n", c.ID)
		}
		return nil
	}

	// Default: print single most recent chat ID
	fmt.Printf("%d\n", latestChatID)
	return nil
}
