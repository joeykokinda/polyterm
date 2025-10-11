package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"polyterm/types"
)

const (
	BaseURL = "https://gamma-api.polymarket.com"
	Timeout = 8 * time.Second
)

func FetchOneEvent(ctx context.Context) (*types.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	url := fmt.Sprintf("%s/events?closed=false&limit=1", BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "polyterm/0.0.1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var events []types.Event
	if err := json.Unmarshal(body, &events); err != nil {
		var eventsResp types.EventsResponse
		if err := json.Unmarshal(body, &eventsResp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		switch {
		case len(eventsResp.Data) > 0:
			events = eventsResp.Data
		case len(eventsResp.Items) > 0:
			events = eventsResp.Items
		default:
			return nil, fmt.Errorf("no events in response")
		}
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no events returned")
	}

	return &events[0], nil
}
