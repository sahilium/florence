#!/usr/bin/env bash

set -e

echo "🛠️ creating folders..."

mkdir -p \
  cmd/daemon \
  internal/{config,event,matrix,ntfy,storage,filters,sinks} \
  sources/{hackernews,lobsters,rss,github,sms}

echo "🧠 writing boilerplate..."

cat > .env <<'EOF'
MATRIX_HOMESERVER=https://matrix-client.matrix.org
MATRIX_ACCESS_TOKEN=changeme
MATRIX_ROOM_ID=!roomid:matrix.org

HN_MIN_SCORE=100
EOF

cat > .gitignore <<'EOF'
.env
*.db
EOF

cat > internal/event/event.go <<'EOF'
package event

import "time"

type Event struct {
	Source    string
	Title     string
	Body      string
	URL       string
	Severity  string
	Timestamp time.Time
	Tags      []string
}
EOF

cat > internal/config/config.go <<'EOF'
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MatrixHomeserver string
	MatrixToken      string
	MatrixRoomID     string
	HNMinScore       int
}

func Load() Config {
	_ = godotenv.Load()

	score, err := strconv.Atoi(getEnv("HN_MIN_SCORE", "100"))
	if err != nil {
		log.Fatal(err)
	}

	return Config{
		MatrixHomeserver: getEnv("MATRIX_HOMESERVER", ""),
		MatrixToken:      getEnv("MATRIX_ACCESS_TOKEN", ""),
		MatrixRoomID:     getEnv("MATRIX_ROOM_ID", ""),
		HNMinScore:       score,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
EOF

cat > internal/sinks/sink.go <<'EOF'
package sinks

import "hn-matrix-bot/internal/event"

type Sink interface {
	Send(event.Event) error
}
EOF

cat > internal/matrix/matrix.go <<'EOF'
package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"hn-matrix-bot/internal/event"
)

type Client struct {
	Homeserver string
	Token      string
	RoomID     string
}

func (c *Client) Send(ev event.Event) error {
	url := fmt.Sprintf(
		"%s/_matrix/client/v3/rooms/%s/send/m.room.message",
		c.Homeserver,
		c.RoomID,
	)

	body := map[string]string{
		"msgtype": "m.text",
		"body": fmt.Sprintf(
			"[%s] %s\n%s",
			ev.Source,
			ev.Title,
			ev.URL,
		),
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
EOF

cat > sources/hackernews/hackernews.go <<'EOF'
package hackernews

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hn-matrix-bot/internal/event"
)

const (
	topStoriesURL = "https://hacker-news.firebaseio.com/v0/topstories.json"
	itemURL       = "https://hacker-news.firebaseio.com/v0/item/%d.json"
)

type Source struct {
	MinScore int
	Seen     map[int]bool
}

type hnItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Score int    `json:"score"`
}

func New(minScore int) *Source {
	return &Source{
		MinScore: minScore,
		Seen:     map[int]bool{},
	}
}

func (s *Source) Fetch() ([]event.Event, error) {
	resp, err := http.Get(topStoriesURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var ids []int
	json.Unmarshal(body, &ids)

	var events []event.Event

	for _, id := range ids[:50] {
		if s.Seen[id] {
			continue
		}

		item := fetchItem(id)
		if item == nil {
			continue
		}

		if item.Score >= s.MinScore {
			events = append(events, event.Event{
				Source:    "hackernews",
				Title:     fmt.Sprintf("🔥 [%d] %s", item.Score, item.Title),
				URL:       item.URL,
				Severity:  "info",
				Timestamp: time.Now(),
			})

			s.Seen[id] = true
		}
	}

	return events, nil
}

func fetchItem(id int) *hnItem {
	url := fmt.Sprintf(itemURL, id)

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var item hnItem
	json.Unmarshal(body, &item)

	return &item
}
EOF

cat > cmd/daemon/main.go <<'EOF'
package main

import (
	"log"
	"time"

	"hn-matrix-bot/internal/config"
	"hn-matrix-bot/internal/matrix"
	"hn-matrix-bot/sources/hackernews"
)

func main() {
	cfg := config.Load()

	matrixClient := &matrix.Client{
		Homeserver: cfg.MatrixHomeserver,
		Token:      cfg.MatrixToken,
		RoomID:     cfg.MatrixRoomID,
	}

	hn := hackernews.New(cfg.HNMinScore)

	log.Println("🚀 daemon started")

	for {
		events, err := hn.Fetch()
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Minute)
			continue
		}

		for _, ev := range events {
			err := matrixClient.Send(ev)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("sent:", ev.Title)
		}

		time.Sleep(5 * time.Minute)
	}
}
EOF
