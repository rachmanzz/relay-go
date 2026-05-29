package relay

import (
	"context"
	"testing"
	"time"
)

func TestRelay_Polling(t *testing.T) {
	r := New()

	t.Run("New Session", func(t *testing.T) {
		ctx := context.Background()
		var callCount int
		
		res := r.Polling(ctx, "", 2, func(lastTimestamp time.Time) any {
			callCount++
			if callCount == 1 {
				return "some-data"
			}
			return nil
		})

		if res["data"] != "some-data" {
			t.Errorf("Expected data 'some-data', got %v", res["data"])
		}
		if res["polling_id"] == "" {
			t.Error("Expected a polling_id to be generated")
		}
	})

	t.Run("Existing Session", func(t *testing.T) {
		ctx := context.Background()
		pollingID := "test-id"
		
		// First call to establish session
		r.Polling(ctx, pollingID, 1, func(ts time.Time) any { return "data1" })

		// Second call should use existing session
		res := r.Polling(ctx, pollingID, 1, func(ts time.Time) any {
			if ts.IsZero() {
				t.Error("Expected non-zero timestamp for existing session")
			}
			return "data2"
		})

		if res["data"] != "data2" {
			t.Errorf("Expected 'data2', got %v", res["data"])
		}
		if res["polling_id"] != pollingID {
			t.Errorf("Expected polling_id %s, got %v", pollingID, res["polling_id"])
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		ctx := context.Background()
		start := time.Now()
		res := r.Polling(ctx, "", 1, func(ts time.Time) any {
			return nil
		})

		duration := time.Since(start)
		if duration < 1*time.Second {
			t.Errorf("Expected timeout after 1s, but took %v", duration)
		}
		if res["data"] != nil {
			t.Errorf("Expected nil data on timeout, got %v", res["data"])
		}
	})

	t.Run("Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(500 * time.Millisecond)
			cancel()
		}()

		res := r.Polling(ctx, "", 5, func(ts time.Time) any {
			return nil
		})

		if res != nil {
			t.Errorf("Expected nil result on context cancellation, got %v", res)
		}
	})
}
