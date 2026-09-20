## Context

See proposal.md for motivation and context. The goal is to provide a clean, standalone Go binary (`tg`) and an associated agent skill to enable both terminal users and AI agents to send messages, upload file attachments (such as generated markdown documents), and statelessly fetch incoming messages.

## Goals / Non-Goals

**Goals:**
- Single standalone Go binary with zero third-party module dependencies, relying purely on the Go standard library.
- Direct integration with Telegram Bot API (endpoints: `/sendMessage`, `/sendDocument`, `/getUpdates`, `/getMe`).
- Support for stdin streaming, file attachments with captions, and parse-mode selection (Markdown, HTML, plain).
- Stateless message fetching with formatted terminal output and machine-readable `--json` output for agent tool consumption.
- Comprehensive Agent Skill (`.agents/skills/telegram/SKILL.md`) accompanied by subpages in `references/` (`sending.md`, `fetching.md`, `setup.md`).

**Non-Goals:**
- MTProto user account client support (no phone number / SMS login).
- Long-running daemon or persistent webhook server.
- Stateful local message database or synchronized offline cache.

## Decisions

### Decision: Pure Go Standard Library Implementation
- **Choice**: Implement the CLI using Go's built-in `net/http`, `mime/multipart`, `encoding/json`, and `flag` packages.
- **Rationale**: Keeps build times negligible, ensures zero supply-chain/dependency vulnerabilities, and produces a compact static binary that is easy to build and distribute.
- **Alternatives Considered**: Cobra/Viper (adds unnecessary weight and external dependencies for a focused CLI).

### Decision: Stateless Polling via `getUpdates`
- **Choice**: The `tg fetch` command queries the Telegram Bot API `getUpdates` endpoint directly without requiring local cursor persistence.
- **Rationale**: Agents and scripts can query recent context without coupling state across execution environments or leaving stale state files on disk. If needed, optional `--offset` and `--limit` flags allow callers to manage update ranges manually.
- **Alternatives Considered**: Storing `last_offset` in `~/.local/state/tg/offset` (adds state drift and multi-consumer race hazards).

### Decision: Agent Skill Architecture with Subpages
- **Choice**: Provide a concise root `SKILL.md` for fast tool resolution and triage, and delegate detailed mechanics to `references/sending.md`, `references/fetching.md`, and `references/setup.md`.
- **Rationale**: Keeps the primary agent context light while giving agents access to deep syntax and formatting details when handling complex file uploads or JSON parsing.

### Decision: Document Upload Mechanics
- **Choice**: Use `sendDocument` with `multipart/form-data` for file uploads, passing filename, optional caption, and parse mode.
- **Rationale**: Telegram's `sendDocument` preserves raw document content (e.g. Markdown files without mangling, code logs, raw text), allowing agents to upload deliverables directly.

## Risks / Trade-offs

- **[Risk] Telegram Markdown Parse Errors**: Sending unescaped special characters in Markdown mode can cause Telegram to reject the message with HTTP 400.
  - *Mitigation*: Fallback to plain text on Markdown parse error if desired, or provide explicit `--parse-mode plain` / auto-escaping guidance in skill documentation.
- **[Risk] Telegram Bot API File Size Limits**: Telegram limits standard Bot API document uploads to 50 MB.
  - *Mitigation*: The CLI inspects file size prior to upload and surfaces a human/agent-readable error if the file exceeds the limit.
- **[Risk] Rate Limiting (HTTP 429)**: Consecutive calls by scripts or agents may trigger Telegram rate limits.
  - *Mitigation*: The HTTP client inspects `Retry-After` headers and retries with backoff for transient 429 responses.
