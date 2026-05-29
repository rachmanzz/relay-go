package relay

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRelay_SSEHandler(t *testing.T) {
	r := New()
	handler := r.SSEHandler()

	t.Run("Headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/sse", nil)
		rr := httptest.NewRecorder()

		// We need to run this in a goroutine because it's a blocking handler
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		
		go handler.ServeHTTP(rr, req.WithContext(ctx))
		
		time.Sleep(50*time.Millisecond)

		if rr.Header().Get("Content-Type") != "text/event-stream" {
			t.Errorf("Expected text/event-stream, got %s", rr.Header().Get("Content-Type"))
		}
	})

	t.Run("Broadcast", func(t *testing.T) {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		client := ts.Client()
		req, _ := http.NewRequest("GET", ts.URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)

		// Read "connected" event
		line, _ := reader.ReadString('\n')
		if !strings.HasPrefix(line, "event: connected") {
			t.Errorf("Expected connected event, got %s", line)
		}

		// Broadcast a test message
		testEvent := "test-event"
		testData := "hello-world"
		
		// Run broadcast in a separate goroutine to ensure client is listening
		go func() {
			time.Sleep(100 * time.Millisecond)
			r.Broadcast(testEvent, testData)
		}()

		// Read the broadcasted event
		// Skip empty lines and "data: ..." lines to find our event
		for i := 0; i < 10; i++ {
			line, _ = reader.ReadString('\n')
			if strings.HasPrefix(line, "event: "+testEvent) {
				line, _ = reader.ReadString('\n')
				if strings.Contains(line, testData) {
					return // Success
				}
			}
		}
		t.Error("Failed to receive broadcasted event")
	})
}
