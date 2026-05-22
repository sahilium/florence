package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sahilium/florence/internal/event"
	"golang.org/x/time/rate"
)

const maxRetries = 3

type Client struct {
	Homeserver string
	Token      string
	Rooms      map[string]string
	limiter    *rate.Limiter
}

type sendResponse struct {
	EventID string `json:"event_id"`
}

type matrixError struct {
	ErrCode      string `json:"errcode"`
	Error        string `json:"error"`
	RetryAfterMs int    `json:"retry_after_ms,omitempty"`
}

func New(homeserver, token string, rooms map[string]string) *Client {
	return &Client{
		Homeserver: homeserver,
		Token:      token,
		Rooms:      rooms,
		limiter:    rate.NewLimiter(rate.Every(time.Second), 1),
	}
}

func (c *Client) Send(ev event.Event) error {
	roomID, ok := c.Rooms[ev.Source]
	if !ok {
		return fmt.Errorf("no room configured for source: %s", ev.Source)
	}

	url := fmt.Sprintf(
		"%s/_matrix/client/v3/rooms/%s/send/m.room.message",
		c.Homeserver,
		roomID,
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

	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if err := c.limiter.Wait(context.Background()); err != nil {
			return fmt.Errorf("rate limit wait: %w", err)
		}

		resp, err := c.doRequest(url, jsonBody)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var sr sendResponse
			if err := json.Unmarshal(respBody, &sr); err != nil {
				return fmt.Errorf("matrix send failed: decode response: %w", err)
			}
			if sr.EventID == "" {
				return fmt.Errorf("matrix send failed: missing event_id in response")
			}
			return nil
		}

		if resp.StatusCode == 429 {
			var mErr matrixError
			if err := json.Unmarshal(respBody, &mErr); err == nil && mErr.RetryAfterMs > 0 {
				lastErr = fmt.Errorf("rate limited: retry_after=%dms", mErr.RetryAfterMs)
				time.Sleep(time.Duration(mErr.RetryAfterMs) * time.Millisecond)
				continue
			}
		}

		return fmt.Errorf(
			"matrix send failed: %d %s",
			resp.StatusCode,
			bytes.TrimSpace(respBody),
		)
	}

	return fmt.Errorf("matrix send failed after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) doRequest(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}

	return resp, nil
}
