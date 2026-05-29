package storage

import (
	"context"
	"time"
)

// Session represents the stateful polling data.
type Session struct {
	LastTimestamp time.Time `json:"last_timestamp"`
	ExpiredAt     time.Time `json:"expired_at"`
}

// Storage defines the interface for session persistence.
type Storage interface {
	// Get retrieves a session by ID.
	Get(ctx context.Context, id string) (*Session, error)
	// Set stores or updates a session.
	Set(ctx context.Context, id string, session *Session, ttl time.Duration) error
	// Delete removes a session.
	Delete(ctx context.Context, id string) error
	// Cleanup removes expired sessions (primarily for in-memory).
	Cleanup(ctx context.Context) error
}
