package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupMockServer() (*httptest.Server, *TelegramClient) {
	mux := http.NewServeMux()

	mux.HandleFunc("/bottester/getMe", func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse[User]{
			OK: true,
			Result: User{
				ID:        123456,
				IsBot:     true,
				FirstName: "VoltaBot",
				Username:  "volta_bot",
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/bottester/sendMessage", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		body, _ := os.ReadFile("/dev/null")
		_ = body
		json.NewDecoder(r.Body).Decode(&payload)

		resp := APIResponse[Message]{
			OK: true,
			Result: Message{
				MessageID: 555,
				Text:      payload["text"].(string),
				Chat: Chat{
					ID: 777,
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/bottester/sendDocument", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseMultipartForm(10 << 20)
		caption := r.FormValue("caption")

		resp := APIResponse[Message]{
			OK: true,
			Result: Message{
				MessageID: 556,
				Caption:   caption,
				Chat: Chat{
					ID: 777,
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/bottester/getUpdates", func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse[[]Update]{
			OK: true,
			Result: []Update{
				{
					UpdateID: 10,
					Message: &Message{
						MessageID: 1001,
						Text:      "Task prompt 1",
						Date:      1700000000,
						From: &User{
							Username: "norty",
						},
						Chat: Chat{
							ID:    777,
							Title: "Test Chat",
							Type:  "group",
						},
					},
				},
				{
					UpdateID: 11,
					Message: &Message{
						MessageID: 1002,
						Text:      "Another chat msg",
						Date:      1700000010,
						From: &User{
							FirstName: "Alice",
						},
						Chat: Chat{
							ID:   888,
							Type: "private",
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	server := httptest.NewServer(mux)
	client := NewTelegramClient("tester")
	client.BaseURL = server.URL + "/bot"
	return server, client
}

func TestRunWhoami(t *testing.T) {
	server, client := setupMockServer()
	defer server.Close()

	ctx := context.Background()
	if err := RunWhoami(ctx, client, false); err != nil {
		t.Fatalf("RunWhoami plain failed: %v", err)
	}
	if err := RunWhoami(ctx, client, true); err != nil {
		t.Fatalf("RunWhoami JSON failed: %v", err)
	}
}

func TestRunChats(t *testing.T) {
	server, client := setupMockServer()
	defer server.Close()

	ctx := context.Background()
	if err := RunChats(ctx, client, false); err != nil {
		t.Fatalf("RunChats plain failed: %v", err)
	}
	if err := RunChats(ctx, client, true); err != nil {
		t.Fatalf("RunChats JSON failed: %v", err)
	}
}

func TestRunSendTextAndStdin(t *testing.T) {
	server, client := setupMockServer()
	defer server.Close()

	ctx := context.Background()

	// 1. Text argument
	err := RunSend(ctx, client, "777", "", "", "markdown", []string{"Hello", "World"}, nil)
	if err != nil {
		t.Fatalf("RunSend args failed: %v", err)
	}

	// 2. Stdin pipe
	stdin := strings.NewReader("Piped text message")
	err = RunSend(ctx, client, "777", "", "", "plain", []string{}, stdin)
	if err != nil {
		t.Fatalf("RunSend stdin failed: %v", err)
	}

	// 3. Error without chat ID
	err = RunSend(ctx, client, "", "", "", "markdown", []string{"test"}, nil)
	if err == nil {
		t.Fatal("expected error with empty chat ID")
	}

	// 4. Error with empty text
	emptyStdin := strings.NewReader("")
	err = RunSend(ctx, client, "777", "", "", "markdown", []string{}, emptyStdin)
	if err == nil {
		t.Fatal("expected error with empty text")
	}
}

func TestRunSendDocument(t *testing.T) {
	server, client := setupMockServer()
	defer server.Close()

	ctx := context.Background()
	tempDir := t.TempDir()
	docPath := filepath.Join(tempDir, "doc.md")
	if err := os.WriteFile(docPath, []byte("# Markdown Document"), 0644); err != nil {
		t.Fatal(err)
	}

	// 1. Valid file
	err := RunSend(ctx, client, "777", docPath, "Caption here", "markdown", nil, nil)
	if err != nil {
		t.Fatalf("RunSend document failed: %v", err)
	}

	// 2. Missing file
	err = RunSend(ctx, client, "777", filepath.Join(tempDir, "missing.md"), "", "markdown", nil, nil)
	if err == nil {
		t.Fatal("expected error with missing document")
	}

	// 3. Stdin as attachment
	stdinDoc := bytes.NewReader([]byte("file from stdin"))
	err = RunSend(ctx, client, "777", "-", "Piped file", "plain", nil, stdinDoc)
	if err != nil {
		t.Fatalf("RunSend document from stdin failed: %v", err)
	}
}

func TestRunFetch(t *testing.T) {
	server, client := setupMockServer()
	defer server.Close()

	ctx := context.Background()

	// 1. Plain fetch
	err := RunFetch(ctx, client, 5, 0, 0, false)
	if err != nil {
		t.Fatalf("RunFetch plain failed: %v", err)
	}

	// 2. JSON fetch
	err = RunFetch(ctx, client, 5, 0, 0, true)
	if err != nil {
		t.Fatalf("RunFetch JSON failed: %v", err)
	}

	// 3. Chat ID filter
	err = RunFetch(ctx, client, 5, 0, 777, true)
	if err != nil {
		t.Fatalf("RunFetch filtered failed: %v", err)
	}
}

func TestRunChatID(t *testing.T) {
	server, client := setupMockServer()
	defer server.Close()

	ctx := context.Background()

	// 1. Single chat ID
	err := RunChatID(ctx, client, false, false)
	if err != nil {
		t.Fatalf("RunChatID default failed: %v", err)
	}

	// 2. All chat IDs
	err = RunChatID(ctx, client, true, false)
	if err != nil {
		t.Fatalf("RunChatID all failed: %v", err)
	}

	// 3. JSON output
	err = RunChatID(ctx, client, false, true)
	if err != nil {
		t.Fatalf("RunChatID JSON failed: %v", err)
	}

	// 4. No updates error
	emptyMux := http.NewServeMux()
	emptyMux.HandleFunc("/botempty/getUpdates", func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse[[]Update]{
			OK:     true,
			Result: []Update{},
		}
		json.NewEncoder(w).Encode(resp)
	})
	emptyServer := httptest.NewServer(emptyMux)
	defer emptyServer.Close()

	emptyClient := NewTelegramClient("empty")
	emptyClient.BaseURL = emptyServer.URL + "/bot"

	err = RunChatID(ctx, emptyClient, false, false)
	if err == nil {
		t.Fatal("expected error when no updates exist")
	}
}
