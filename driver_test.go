package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDriver_Fallback(t *testing.T) {
	// Mock SSE Server that fails after initial connection
	sseServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "event: connected\ndata: test-session-id\n\n")
		w.(http.Flusher).Flush()
		
		time.Sleep(100 * time.Millisecond)
		// Simulate failure by closing connection
	}))
	defer sseServer.Close()

	// Mock Polling Server
	pollingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		pollingID := req.URL.Query().Get("polling_id")
		res := map[string]any{
			"polling_id": pollingID,
			"data":       "polling-data",
			"timestamp":  time.Now(),
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer pollingServer.Close()

	driver := NewDriver(sseServer.URL, pollingServer.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	driver.Start(ctx)

	// Wait for event from polling after SSE fails
	select {
	case ev := <-driver.Events():
		if ev.Data != "polling-data" {
			t.Errorf("Expected polling-data, got %v", ev.Data)
		}
	case <-time.After(3 * time.Second):
		t.Error("Timed out waiting for fallback event")
	}
}
