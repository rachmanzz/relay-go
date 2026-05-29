package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (s *RedisStorage) Get(ctx context.Context, id string) (*Session, error) {
	data, err := s.client.Get(ctx, "relay:session:"+id).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *RedisStorage) Set(ctx context.Context, id string, session *Session, ttl time.Duration) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, "relay:session:"+id, data, ttl).Err()
}

func (s *RedisStorage) Delete(ctx context.Context, id string) error {
	return s.client.Del(ctx, "relay:session:"+id).Err()
}

func (s *RedisStorage) Cleanup(ctx context.Context) error {
	// Redis handles expiration via TTL automatically
	return nil
}
