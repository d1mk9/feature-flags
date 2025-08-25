package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"feature-flags/pkg/repository"

	lru "github.com/hashicorp/golang-lru/v2"
)

type Flags interface {
	GetValue(ctx context.Context, key string) (FeatureValue, error)
	SetValue(ctx context.Context, key string, value FeatureValue) error
}

type cachedItem struct {
	raw      json.RawMessage
	cachedAt time.Time
	ttl      time.Duration
}

type FeatureService struct {
	repo  repository.Repository
	cache *lru.Cache[string, cachedItem]
	ttl   time.Duration
}

func NewFeatureService(repo repository.Repository, cacheSize int, ttlMinutes int) (*FeatureService, error) {
	c, err := lru.New[string, cachedItem](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("create LRU cache: %w", err)
	}

	return &FeatureService{
		repo:  repo,
		cache: c,
		ttl:   time.Duration(ttlMinutes) * time.Minute,
	}, nil
}

func (s *FeatureService) getFromCache(key string) (json.RawMessage, bool) {
	if ci, ok := s.cache.Get(key); ok && time.Since(ci.cachedAt) < ci.ttl {
		return ci.raw, true
	}

	s.cache.Remove(key)
	return nil, false
}

func (s *FeatureService) putToCache(key string, raw json.RawMessage) {
	s.cache.Add(key, cachedItem{raw: raw, cachedAt: time.Now(), ttl: s.ttl})
}

func (s *FeatureService) GetValue(ctx context.Context, key string) (FeatureValue, error) {
	if raw, ok := s.getFromCache(key); ok {
		var v FeatureValue
		if err := json.Unmarshal(raw, &v); err != nil {
			return FeatureValue{}, err
		}
		return v, nil
	}

	raw, err := s.repo.GetValueByKey(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return FeatureValue{}, ErrNotFound
		}
		return FeatureValue{}, fmt.Errorf("get value: %w", err)
	}
	s.putToCache(key, raw)

	var v FeatureValue
	if err := json.Unmarshal(raw, &v); err != nil {
		return FeatureValue{}, err
	}
	return v, nil
}

func (s *FeatureService) SetValue(ctx context.Context, key string, value FeatureValue) error {
	if !value.Validate() {
		return ErrUnsupportedType
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode value: %w", err)
	}
	if err := s.repo.SetValueByKey(ctx, key, raw); err != nil {
		return fmt.Errorf("set value: %w", err)
	}
	s.cache.Remove(key)
	return nil
}
