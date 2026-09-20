package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds runtime configuration options for the Telegram CLI.
type Config struct {
	BotToken string `json:"bot_token"`
	ChatID   string `json:"chat_id"`
}

// DefaultConfigPath returns ~/.config/tg/config.json.
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "tg", "config.json")
}

// LoadConfigFile reads a JSON config file from the specified path.
func LoadConfigFile(path string) (*Config, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return &cfg, nil
}

// ResolveConfig resolves configuration with precedence:
// 1. Explicit CLI flag
// 2. Environment variables (TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID)
// 3. Config file (custom path if provided, else DefaultConfigPath())
func ResolveConfig(flagToken, flagChatID, customConfigPath string) (*Config, error) {
	cfg := &Config{}

	configPath := customConfigPath
	if configPath == "" {
		configPath = os.Getenv("TG_CONFIG")
	}
	if configPath == "" {
		configPath = DefaultConfigPath()
	}

	fileCfg, err := LoadConfigFile(configPath)
	if err != nil {
		return nil, err
	}

	if fileCfg != nil {
		cfg.BotToken = fileCfg.BotToken
		cfg.ChatID = fileCfg.ChatID
	}

	if envToken := os.Getenv("TELEGRAM_BOT_TOKEN"); envToken != "" {
		cfg.BotToken = envToken
	}
	if envChatID := os.Getenv("TELEGRAM_CHAT_ID"); envChatID != "" {
		cfg.ChatID = envChatID
	}

	if flagToken != "" {
		cfg.BotToken = flagToken
	}
	if flagChatID != "" {
		cfg.ChatID = flagChatID
	}

	return cfg, nil
}
