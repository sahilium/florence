package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sahilium/florence/internal/event"
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
