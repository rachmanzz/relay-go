package storage

import (
	"context"
	"sync"
	"time"
)

type InMemoryStorage struct {
	sessions sync.Map
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

func (s *InMemoryStorage) Get(ctx context.Context, id string) (*Session, error) {
	val, ok := s.sessions.Load(id)
	if !ok {
		return nil, nil
	}
	session := val.(*Session)
	if time.Now().After(session.ExpiredAt) {
		s.sessions.Delete(id)
		return nil, nil
	}
	return session, nil
}

func (s *InMemoryStorage) Set(ctx context.Context, id string, session *Session, ttl time.Duration) error {
	s.sessions.Store(id, session)
	return nil
}

func (s *InMemoryStorage) Delete(ctx context.Context, id string) error {
	s.sessions.Delete(id)
	return nil
}

func (s *InMemoryStorage) Cleanup(ctx context.Context) error {
	now := time.Now()
	s.sessions.Range(func(key, value any) bool {
		session := value.(*Session)
		if now.After(session.ExpiredAt) {
			s.sessions.Delete(key)
		}
		return true
	})
	return nil
}
