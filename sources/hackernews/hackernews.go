package hackernews

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sahilium/florence/internal/event"
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
				ID:        fmt.Sprintf("hn_%d", item.ID),
				Source:    "hackernews",
				Title:     fmt.Sprintf("[%d] %s", item.Score, item.Title),
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
