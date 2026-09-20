# GitHub Copilot Instructions

## Project Context
This repository contains `tg`, a lightweight, dependency-free CLI tool and agent skill written in Go for Telegram messaging, document uploading, and stateless update polling.

## Tool Execution Guidelines
When the user asks to send notifications, dispatch deliverables, or read messages via Telegram:

1. **Sending Files & Markdown Reports**:
   - Prefer uploading markdown documents or logs directly rather than sending large text walls:
     ```bash
     tg send -f path/to/report.md -c "Release Notes"
     ```
2. **Sending Status Updates**:
   ```bash
   tg send "Deployment finished with exit status 0"
   ```
3. **Retrieving Chat ID**:
   ```bash
   tg chat-id --token "<bot-token>"
   ```
4. **Fetching Incoming Messages**:
   - Always use `--json` when machine parsing is required:
     ```bash
     tg fetch --limit 5 --json
     ```

## Project Standards
- Language: Go 1.19+ (Standard library only; zero external dependencies).
- Build: `make build` (compiles to `bin/tg`).
- Tests: `make test`.
- Reference skill definitions: [`.agents/skills/telegram/SKILL.md`](../.agents/skills/telegram/SKILL.md).
