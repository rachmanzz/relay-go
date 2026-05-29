package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rachmanzz/relay-go/internal/session/storage"
)

type Manager struct {
	storage storage.Storage
	ttl     time.Duration
}

func NewManager(s storage.Storage, ttl time.Duration) *Manager {
	m := &Manager{
		storage: s,
		ttl:     ttl,
	}

	// Start background cleanup for in-memory storage
	go m.startCleanup()

	return m
}

func (m *Manager) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		_ = m.storage.Cleanup(context.Background())
	}
}

func (m *Manager) HandlePolling(ctx context.Context, pollingID string, timeout int, fn func(lastTimestamp time.Time) any) map[string]any {
	var session *storage.Session
	var err error

	if pollingID != "" {
		session, err = m.storage.Get(ctx, pollingID)
		if err != nil {
			// Log error if needed, but fallback to new session
			session = nil
		}
	}

	if session == nil {
		if pollingID == "" {
			pollingID = uuid.New().String()
		}
		session = &storage.Session{
			LastTimestamp: time.Now().AddDate(0, 0, -5),
		}
	}

	session.ExpiredAt = time.Now().Add(m.ttl)
	_ = m.storage.Set(ctx, pollingID, session, m.ttl)

	startTime := time.Now()
	timeoutDuration := time.Duration(timeout) * time.Second

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			data := fn(session.LastTimestamp)
			if data != nil {
				session.LastTimestamp = time.Now()
				session.ExpiredAt = time.Now().Add(m.ttl)
				_ = m.storage.Set(ctx, pollingID, session, m.ttl)

				return map[string]any{
					"polling_id": pollingID,
					"data":       data,
					"timestamp":  session.LastTimestamp,
				}
			}

			if time.Since(startTime) >= timeoutDuration {
				return map[string]any{
					"polling_id": pollingID,
					"data":       nil,
					"timestamp":  session.LastTimestamp,
				}
			}

			time.Sleep(1 * time.Second)
		}
	}
}
