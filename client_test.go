package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientEndpoints(t *testing.T) {
	mux := http.NewServeMux()

	// Mock getMe
	mux.HandleFunc("/bottester/getMe", func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse[User]{
			OK: true,
			Result: User{
				ID:        123456,
				IsBot:     true,
				FirstName: "TestBot",
				Username:  "test_bot",
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// Mock sendMessage
	mux.HandleFunc("/bottester/sendMessage", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &payload)

		resp := APIResponse[Message]{
			OK: true,
			Result: Message{
				MessageID: 101,
				Text:      payload["text"].(string),
				Chat: Chat{
					ID: 999999,
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// Mock sendDocument
	mux.HandleFunc("/bottester/sendDocument", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		file, header, _ := r.FormFile("document")
		defer file.Close()
		content, _ := io.ReadAll(file)

		resp := APIResponse[Message]{
			OK: true,
			Result: Message{
				MessageID: 102,
				Caption:   r.FormValue("caption"),
				Text:      string(content),
				Chat: Chat{
					ID: 999999,
				},
			},
		}
		_ = header
		json.NewEncoder(w).Encode(resp)
	})

	// Mock getUpdates
	mux.HandleFunc("/bottester/getUpdates", func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse[[]Update]{
			OK: true,
			Result: []Update{
				{
					UpdateID: 1,
					Message: &Message{
						MessageID: 201,
						Text:      "Hello from user",
						From: &User{
							ID:        888,
							Username:  "tester",
							FirstName: "Test",
						},
						Chat: Chat{
							ID:    999999,
							Title: "Test Group",
						},
						Date: 1700000000,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewTelegramClient("tester")
	client.BaseURL = server.URL + "/bot"

	ctx := context.Background()

	// 1. Test GetMe
	user, err := client.GetMe(ctx)
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}
	if user.Username != "test_bot" {
		t.Errorf("expected test_bot, got %s", user.Username)
	}

	// 2. Test SendMessage
	msg, err := client.SendMessage(ctx, "999999", "Test Hello", "markdown")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	if msg.Text != "Test Hello" || msg.MessageID != 101 {
		t.Errorf("unexpected msg output: %+v", msg)
	}

	// 3. Test SendDocument
	fileContent := []byte("# Markdown Report\nAll tests passed.")
	docMsg, err := client.SendDocument(ctx, "999999", "report.md", bytes.NewReader(fileContent), "My Caption", "markdown")
	if err != nil {
		t.Fatalf("SendDocument failed: %v", err)
	}
	if docMsg.Caption != "My Caption" || docMsg.MessageID != 102 {
		t.Errorf("unexpected docMsg output: %+v", docMsg)
	}

	// 4. Test GetUpdates
	updates, err := client.GetUpdates(ctx, 0, 5)
	if err != nil {
		t.Fatalf("GetUpdates failed: %v", err)
	}
	if len(updates) != 1 || updates[0].EffectiveMessage().Text != "Hello from user" {
		t.Errorf("unexpected updates: %+v", updates)
	}
}
