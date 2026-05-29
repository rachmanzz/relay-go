package transport

import (
	"fmt"
	"net/http"
)

// SSEHandler handles Server-Sent Events connections.
type SSEHandler struct {
}

func NewSSEHandler() *SSEHandler {
	return &SSEHandler{}
}

// ServeHTTP implements the http.Handler interface.
func (h *SSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Flush the headers immediately
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	flusher.Flush()

	// Keep the connection open
	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

// WriteEvent writes an SSE formatted event to the response writer.
func WriteEvent(w http.ResponseWriter, event string, data string) error {
	if event != "" {
		if _, err := fmt.Fprintf(w, "event: %s\n", event); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
		return err
	}
	
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}
