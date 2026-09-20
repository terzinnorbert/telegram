package main

import (
	"time"
)

// APIResponse wraps standard Telegram Bot API response format.
type APIResponse[T any] struct {
	OK          bool                `json:"ok"`
	Result      T                   `json:"result"`
	ErrorCode   int                 `json:"error_code,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

// ResponseParameters contains details about why a request was unsuccessful.
type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      int   `json:"retry_after,omitempty"`
}

// User represents a Telegram user or bot.
type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// Chat represents a Telegram chat.
type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// Message represents a Telegram message.
type Message struct {
	MessageID int64  `json:"message_id"`
	From      *User  `json:"from,omitempty"`
	Date      int64  `json:"date"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text,omitempty"`
	Caption   string `json:"caption,omitempty"`
}

// Update represents an incoming update from getUpdates.
type Update struct {
	UpdateID    int64    `json:"update_id"`
	Message     *Message `json:"message,omitempty"`
	ChannelPost *Message `json:"channel_post,omitempty"`
}

// JSONMessage is the structured JSON model returned to agents and scripts.
type JSONMessage struct {
	UpdateID  int64  `json:"update_id"`
	MessageID int64  `json:"message_id"`
	ChatID    int64  `json:"chat_id"`
	ChatTitle string `json:"chat_title,omitempty"`
	Sender    string `json:"sender"`
	Date      string `json:"date"`
	Text      string `json:"text"`
}

// ChatSummary represents a discovered chat in `tg chats`.
type ChatSummary struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title,omitempty"`
	Username string `json:"username,omitempty"`
}

// EffectiveMessage returns Message or ChannelPost if available.
func (u *Update) EffectiveMessage() *Message {
	if u.Message != nil {
		return u.Message
	}
	return u.ChannelPost
}

// ToJSONMessage converts an Update into JSONMessage representation.
func (u *Update) ToJSONMessage() *JSONMessage {
	msg := u.EffectiveMessage()
	if msg == nil {
		return nil
	}

	sender := ""
	if msg.From != nil {
		if msg.From.Username != "" {
			sender = "@" + msg.From.Username
		} else {
			sender = msg.From.FirstName
			if msg.From.LastName != "" {
				sender += " " + msg.From.LastName
			}
		}
	} else if msg.Chat.Title != "" {
		sender = msg.Chat.Title
	}

	text := msg.Text
	if text == "" {
		text = msg.Caption
	}

	return &JSONMessage{
		UpdateID:  u.UpdateID,
		MessageID: msg.MessageID,
		ChatID:    msg.Chat.ID,
		ChatTitle: msg.Chat.Title,
		Sender:    sender,
		Date:      time.Unix(msg.Date, 0).UTC().Format(time.RFC3339),
		Text:      text,
	}
}
