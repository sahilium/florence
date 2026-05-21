## Florence

[![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Matrix](https://img.shields.io/badge/matrix-enabled-000000?logo=matrix)](https://matrix.org/)
[![SQLite](https://img.shields.io/badge/storage-sqlite-003B57?logo=sqlite)](https://sqlite.org/)
[![Status](https://img.shields.io/badge/status-active-success.svg)]()

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
    Source    string
    Title     string
    Body      string
    URL       string
    Severity  string
    Timestamp time.Time
    Tags      []string
}
```

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
- Matrix account
- Matrix room
- Matrix access token

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

Create a ".env" file:

MATRIX_HOMESERVER=https://matrix-client.matrix.org
MATRIX_ACCESS_TOKEN=your_access_token
MATRIX_ROOM_ID=!roomid:matrix.org

HN_MIN_SCORE=100

---

### Run

```sh
go run ./cmd/daemon
```

---

## Matrix Integration

The daemon can publish notifications directly into Matrix rooms using the Matrix Client API.

Current implementation supports:

- Plain text notifications
- Room-based delivery
- Access token authentication

## Planned improvements:

- Rich formatting
- Markdown rendering
- Threaded notifications
- Multiple room routing

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

- SQLite persistence
- Deduplication
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