package relay

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Driver represents the client-side intelligent switcher.
type Driver struct {
	url        string
	pollingURL string
	client     *http.Client
	
	mu         sync.RWMutex
	pollingID  string
	lastTS     time.Time
	
	events     chan Event
	stop       chan struct{}
}

type Event struct {
	Name string
	Data any
}

func NewDriver(url string, pollingURL string) *Driver {
	return &Driver{
		url:        url,
		pollingURL: pollingURL,
		client:     &http.Client{},
		events:     make(chan Event, 100),
		stop:       make(chan struct{}),
	}
}

func (d *Driver) Events() <-chan Event {
	return d.events
}

func (d *Driver) Start(ctx context.Context) {
	go d.run(ctx)
}

func (d *Driver) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-d.stop:
			return
		default:
			// Try SSE first
			err := d.connectSSE(ctx)
			if err != nil {
				// Fallback to Polling
				d.runPolling(ctx)
			}
		}
	}
}

func (d *Driver) connectSSE(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", d.url, nil)
	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				return err // Connection lost, trigger fallback/reconnect
			}

			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "event:") {
				eventName := strings.TrimSpace(strings.TrimPrefix(line, "event:"))
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				data := strings.TrimSpace(strings.TrimPrefix(dataLine, "data:"))
				
				if eventName == "connected" {
					d.mu.Lock()
					d.pollingID = data
					d.mu.Unlock()
					continue
				}

				d.events <- Event{Name: eventName, Data: data}
			}
		}
	}
}

func (d *Driver) runPolling(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			d.mu.RLock()
			pid := d.pollingID
			d.mu.RUnlock()

			url := fmt.Sprintf("%s?polling_id=%s", d.pollingURL, pid)
			req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
			resp, err := d.client.Do(req)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			defer resp.Body.Close()

			var result map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				time.Sleep(2 * time.Second)
				continue
			}

			if result["polling_id"] != nil {
				d.mu.Lock()
				d.pollingID = result["polling_id"].(string)
				d.mu.Unlock()
			}

			if data := result["data"]; data != nil {
				d.events <- Event{Name: "message", Data: data}
			}

			// If we should try to upgrade back to SSE, we could add logic here.
			// For now, keep polling until the outer loop restarts (e.g. on error).
			time.Sleep(1 * time.Second)
		}
	}
}
