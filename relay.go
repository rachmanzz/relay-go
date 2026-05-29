package relay

import (
	"context"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rachmanzz/relay-go/internal/event"
	"github.com/rachmanzz/relay-go/internal/session"
	"github.com/rachmanzz/relay-go/internal/session/storage"
	"github.com/rachmanzz/relay-go/internal/transport"
	"github.com/google/uuid"
)

type Options struct {
	RedisClient *redis.Client
	SessionTTL  time.Duration
}

var defaultOptions = Options{
	SessionTTL: 5 * time.Minute,
}

type Relay struct {
	sessionManager *session.Manager
	broadcaster    *event.Broadcaster
}

func New(opts ...func(*Options)) *Relay {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	var store storage.Storage
	if o.RedisClient != nil {
		store = storage.NewRedisStorage(o.RedisClient)
	} else {
		store = storage.NewInMemoryStorage()
	}

	return &Relay{
		sessionManager: session.NewManager(store, o.SessionTTL),
		broadcaster:    event.NewBroadcaster(),
	}
}

func WithRedis(client *redis.Client) func(*Options) {
	return func(o *Options) {
		o.RedisClient = client
	}
}

func WithTTL(ttl time.Duration) func(*Options) {
	return func(o *Options) {
		o.SessionTTL = ttl
	}
}

// Polling is the framework-agnostic helper to handle stateful long polling.
func (r *Relay) Polling(ctx context.Context, pollingID string, timeout int, fn func(lastTimestamp time.Time) any) map[string]any {
	return r.sessionManager.HandlePolling(ctx, pollingID, timeout, fn)
}

// SSEHandler returns a standard http.Handler for SSE connections.
func (r *Relay) SSEHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		id := uuid.New().String()
		msgChan := r.broadcaster.Subscribe(id)
		defer r.broadcaster.Unsubscribe(id)

		// Send initial connected event
		_ = transport.WriteEvent(w, "connected", id)
		flusher.Flush()

		ctx := req.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgChan:
				if !ok {
					return
				}
				_ = transport.WriteEvent(w, msg.Event, msg.Data)
				flusher.Flush()
			}
		}
	})
}

// Broadcast sends an event to all connected SSE clients.
func (r *Relay) Broadcast(evt string, data string) {
	r.broadcaster.Broadcast(event.Message{Event: evt, Data: data})
}

// SendTo sends an event to a specific client ID.
func (r *Relay) SendTo(clientID string, evt string, data string) {
	r.broadcaster.SendTo(clientID, event.Message{Event: evt, Data: data})
}

