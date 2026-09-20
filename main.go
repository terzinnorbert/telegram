package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
)

const version = "1.0.0"

func printUsage() {
	fmt.Fprintf(os.Stderr, `tg - Fast, lightweight Telegram CLI for users and AI agents (v%s)

USAGE:
  tg <command> [arguments...]

COMMANDS:
  send        Send a text message or upload a file
  fetch       Statelessly retrieve recent incoming messages
  chats       List active chats discovered from recent updates
  chat-id     Retrieve active chat ID directly (alias: get-chat-id)
  whoami      Display bot username and account information
  version     Show version information

GLOBAL FLAGS:
  --token <token>       Telegram Bot token (or TELEGRAM_BOT_TOKEN env var)
  --chat-id <id>        Default chat ID (or TELEGRAM_CHAT_ID env var)
  --config <path>       Path to config file (default: ~/.config/tg/config.json)

COMMAND HELP:
  Run 'tg <command> --help' for details on specific command options.

EXAMPLES:
  tg chat-id --token <token>
  tg send "Hello from CLI"
  tg send -f report.md -c "Sprint Report"
  git status | tg send
  tg fetch --limit 5 --json
  tg chats
`, version)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "help", "-h", "--help":
		printUsage()
		return
	case "version", "-v", "--version":
		fmt.Printf("tg version %s\n", version)
		return
	}

	ctx := context.Background()

	switch cmd {
	case "whoami":
		fs := flag.NewFlagSet("whoami", flag.ExitOnError)
		flagToken := fs.String("token", "", "Telegram bot token")
		flagConfig := fs.String("config", "", "Config file path")
		flagJSON := fs.Bool("json", false, "Output as JSON")
		_ = fs.Parse(args)

		cfg, err := ResolveConfig(*flagToken, "", *flagConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		if cfg.BotToken == "" {
			fmt.Fprintf(os.Stderr, "Error: Bot token required. Set TELEGRAM_BOT_TOKEN or pass --token\n")
			os.Exit(1)
		}

		client := NewTelegramClient(cfg.BotToken)
		if err := RunWhoami(ctx, client, *flagJSON); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "chats":
		fs := flag.NewFlagSet("chats", flag.ExitOnError)
		flagToken := fs.String("token", "", "Telegram bot token")
		flagConfig := fs.String("config", "", "Config file path")
		flagJSON := fs.Bool("json", false, "Output as JSON")
		_ = fs.Parse(args)

		cfg, err := ResolveConfig(*flagToken, "", *flagConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		if cfg.BotToken == "" {
			fmt.Fprintf(os.Stderr, "Error: Bot token required. Set TELEGRAM_BOT_TOKEN or pass --token\n")
			os.Exit(1)
		}

		client := NewTelegramClient(cfg.BotToken)
		if err := RunChats(ctx, client, *flagJSON); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "chat-id", "get-chat-id":
		fs := flag.NewFlagSet("chat-id", flag.ExitOnError)
		flagToken := fs.String("token", "", "Telegram bot token")
		flagConfig := fs.String("config", "", "Config file path")
		flagAll := fs.Bool("all", false, "Output all discovered unique chat IDs")
		flagJSON := fs.Bool("json", false, "Output discovered chats as JSON")
		_ = fs.Parse(args)

		token := *flagToken
		if token == "" && len(fs.Args()) > 0 {
			token = fs.Args()[0]
		}

		cfg, err := ResolveConfig(token, "", *flagConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		if cfg.BotToken == "" {
			fmt.Fprintf(os.Stderr, "Error: Bot token required. Pass --token, supply token argument, or set TELEGRAM_BOT_TOKEN\n")
			os.Exit(1)
		}

		client := NewTelegramClient(cfg.BotToken)
		if err := RunChatID(ctx, client, *flagAll, *flagJSON); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "send":
		fs := flag.NewFlagSet("send", flag.ExitOnError)
		flagToken := fs.String("token", "", "Telegram bot token")
		flagChatID := fs.String("chat-id", "", "Target chat ID")
		flagConfig := fs.String("config", "", "Config file path")
		flagFile := fs.String("f", "", "File attachment path (or '-' for stdin)")
		flagFileLong := fs.String("file", "", "File attachment path (or '-' for stdin)")
		flagCaption := fs.String("c", "", "Caption for attachment")
		flagCaptionLong := fs.String("caption", "", "Caption for attachment")
		flagParseMode := fs.String("parse-mode", "markdown", "Parse mode: markdown, html, plain")
		_ = fs.Parse(args)

		targetFile := *flagFile
		if targetFile == "" {
			targetFile = *flagFileLong
		}
		caption := *flagCaption
		if caption == "" {
			caption = *flagCaptionLong
		}

		cfg, err := ResolveConfig(*flagToken, *flagChatID, *flagConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		if cfg.BotToken == "" {
			fmt.Fprintf(os.Stderr, "Error: Bot token required. Set TELEGRAM_BOT_TOKEN or pass --token\n")
			os.Exit(1)
		}
		if cfg.ChatID == "" {
			fmt.Fprintf(os.Stderr, "Error: Chat ID required. Set TELEGRAM_CHAT_ID or pass --chat-id\n")
			os.Exit(1)
		}

		client := NewTelegramClient(cfg.BotToken)
		if err := RunSend(ctx, client, cfg.ChatID, targetFile, caption, *flagParseMode, fs.Args(), os.Stdin); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "fetch":
		fs := flag.NewFlagSet("fetch", flag.ExitOnError)
		flagToken := fs.String("token", "", "Telegram bot token")
		flagConfig := fs.String("config", "", "Config file path")
		flagChatID := fs.String("chat-id", "", "Filter messages by chat ID")
		flagLimit := fs.Int("limit", 10, "Maximum number of messages to retrieve")
		flagOffset := fs.Int64("offset", 0, "Update ID offset")
		flagJSON := fs.Bool("json", false, "Output as structured JSON")
		_ = fs.Parse(args)

		cfg, err := ResolveConfig(*flagToken, "", *flagConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		if cfg.BotToken == "" {
			fmt.Fprintf(os.Stderr, "Error: Bot token required. Set TELEGRAM_BOT_TOKEN or pass --token\n")
			os.Exit(1)
		}

		var chatFilter int64
		if *flagChatID != "" {
			parsed, err := strconv.ParseInt(*flagChatID, 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid --chat-id '%s': %v\n", *flagChatID, err)
				os.Exit(1)
			}
			chatFilter = parsed
		}

		client := NewTelegramClient(cfg.BotToken)
		if err := RunFetch(ctx, client, *flagLimit, *flagOffset, chatFilter, *flagJSON); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nRun 'tg --help' for usage.\n", cmd)
		os.Exit(1)
	}
}
