package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPrecedence(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")
	configJSON := `{"bot_token": "token_from_file", "chat_id": "chat_from_file"}`
	if err := os.WriteFile(configFile, []byte(configJSON), 0644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Case 1: Config file only
	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	os.Unsetenv("TELEGRAM_CHAT_ID")
	cfg, err := ResolveConfig("", "", configFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BotToken != "token_from_file" || cfg.ChatID != "chat_from_file" {
		t.Errorf("expected config from file, got %+v", cfg)
	}

	// Case 2: Environment variable overrides config file
	os.Setenv("TELEGRAM_BOT_TOKEN", "token_from_env")
	os.Setenv("TELEGRAM_CHAT_ID", "chat_from_env")
	cfg, err = ResolveConfig("", "", configFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BotToken != "token_from_env" || cfg.ChatID != "chat_from_env" {
		t.Errorf("expected env to override file, got %+v", cfg)
	}

	// Case 3: Flag overrides environment variable and config file
	cfg, err = ResolveConfig("token_from_flag", "chat_from_flag", configFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BotToken != "token_from_flag" || cfg.ChatID != "chat_from_flag" {
		t.Errorf("expected flag to override env and file, got %+v", cfg)
	}

	// Cleanup
	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	os.Unsetenv("TELEGRAM_CHAT_ID")
}
