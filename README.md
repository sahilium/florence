## Florence

[![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8?logo=go&style=for-the-badge)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)](LICENSE)
[![Matrix](https://img.shields.io/badge/matrix-enabled-000000?logo=matrix&style=for-the-badge)](https://matrix.org/)
[![SQLite](https://img.shields.io/badge/storage-sqlite-003B57?logo=sqlite&style=for-the-badge)](https://sqlite.org/)
[![Status](https://img.shields.io/badge/status-active-success.svg?style=for-the-badge)]()

A lightweight event ingestion and notification platform written in Go.

Florence aggregates events from multiple sources such as Hacker News, Lobsters, RSS feeds, GitHub, SMS, and other services, normalizes them into a unified internal event model, and routes them to configurable notification sinks like Matrix, ntfy, Discord, and more.

---

## Features

- Modular source adapters
- Unified internal event model
- Multiple notification sinks
- Matrix integration
- Configurable filtering and routing
- Environment variable based configuration
- Lightweight and self-hostable
- SQLite-friendly architecture
- Minimal dependencies
- Easy extension for new event sources

---

## Event Model

All incoming events are normalized into a common structure.

```go
type Event struct {
    ID        string
    Source    string
    Title     string
    Body      string
    URL       string
    Severity  string
    Timestamp time.Time
    Tags      []string
}

Events carry deterministic IDs (e.g. `hn_12345`) for persistent deduplication.
This abstraction allows any source to route events through any notification sink.

---

## Supported Sources

| Source       | Status      |
|--------------|-------------|
| Hacker News  | Implemented |
| Lobsters     | Planned     |
| RSS          | Planned     |
| GitHub       | Planned     |
| SMS          | Planned     |

---

## Supported Sinks

| Sink     | Status      |
|----------|-------------|
| Matrix   | Implemented |
| ntfy     | Planned     |
| Discord  | Planned     |
| Email    | Planned     |

---

## Getting Started

### Prerequisites

- Go 1.24+
- Matrix account and access token
- Turso database (or libSQL-compatible server)

---

### Clone Repository

```sh
git clone https://github.com/sahilium/florence.git
cd florence
```

---

### Install Dependencies

```sh
go mod tidy
```

---

### Configure Environment Variables

Create a ".env" file, see [.env.example](.env.example) for an example.

---

### Run

```sh
go run ./cmd/daemon
```

---

## Matrix Integration

The daemon can publish notifications directly into Matrix rooms using the Matrix Client API.

Events are routed dynamically: each event source publishes to its own Matrix room based on the `ROOM_*` environment variables.

Current implementation supports:

- Plain text notifications
- Per-source room routing
- Access token authentication

## Planned improvements:

- Rich formatting
- Markdown rendering
- Threaded notifications

---

## Database

Florence uses Turso/libSQL for persistent storage.

### Migrations

Database migrations are managed by [Goose](https://github.com/pressly/goose) and stored in `migrations/`. Migrations run automatically on daemon startup.

### Tables

- `seen_events` — deterministic deduplication (keyed by event ID).
- `source_state` — per-source polling cursors and state.
- `events` — normalized event history for debugging.

### Deduplication Flow

```
source → normalize event → check seen_events
  → if seen: skip
  → else: store event → send to sink → mark as seen
```

---

## Design Goals

- Small deployable footprint
- Minimal infrastructure requirements
- Event-driven architecture
- Modular source/sink abstractions
- Easy self-hosting
- Simple operational model
- Extensibility without framework lock-in

---

## Future Plans

- Rule engine
- AI summarization
- Web dashboard
- Cloudflare Workers deployment
- Multi-user support
- Webhook ingestion
- Event replay
- Metrics and observability

---

## License

MIT
