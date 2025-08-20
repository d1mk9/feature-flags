package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"feature-flags/pkg/repository"

	lru "github.com/hashicorp/golang-lru/v2"
)

type Vars interface {
	GetValue(ctx context.Context, key string) (json.RawMessage, error)
	SetValue(ctx context.Context, key string, value json.RawMessage) error
}

type cachedItem struct {
	val      json.RawMessage
	cachedAt time.Time
	ttl      time.Duration
}

type FeatureService struct {
	repo    repository.Repository
	advRepo repository.AdvancedRepository
	cache   *lru.Cache[string, cachedItem]
	ttl     time.Duration
}

func NewFeatureService(repo repository.Repository, cacheSize int, ttlMinutes int) (*FeatureService, error) {
	c, err := lru.New[string, cachedItem](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("create LRU cache: %w", err)
	}
	var adv repository.AdvancedRepository
	if a, ok := repo.(repository.AdvancedRepository); ok {
		adv = a
	}
	return &FeatureService{
		repo:    repo,
		advRepo: adv,
		cache:   c,
		ttl:     time.Duration(ttlMinutes) * time.Minute,
	}, nil
}

func (s *FeatureService) GetValue(ctx context.Context, key string) (json.RawMessage, error) {
	if ci, ok := s.cache.Get(key); ok {
		if time.Since(ci.cachedAt) < ci.ttl {
			return ci.val, nil
		}
		s.cache.Remove(key)
	}

	value, err := s.repo.GetValue(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("service.GetValue key=%q: %w", key, err)
	}

	s.cache.Add(key, cachedItem{
		val:      value,
		cachedAt: time.Now(),
		ttl:      s.ttl,
	})
	return value, nil
}

func (s *FeatureService) SetValue(ctx context.Context, key string, value json.RawMessage) error {
	if err := s.repo.UpsertValue(ctx, key, value); err != nil {
		return fmt.Errorf("service.SetValue key=%q: %w", key, err)
	}
	s.cache.Remove(key)
	return nil
}

var ErrNotImplemented = repository.ErrNotImplemented

func (s *FeatureService) ListFeatures(ctx context.Context, f repository.FeatureListFilters) ([]repository.FeatureRow, error) {
	if s.advRepo == nil {
		return nil, ErrNotImplemented
	}
	return s.advRepo.ListFeatures(ctx, f)
}

func (s *FeatureService) CountFeatures(ctx context.Context, f repository.FeatureListFilters) (int64, error) {
	if s.advRepo == nil {
		return 0, ErrNotImplemented
	}
	return s.advRepo.CountFeatures(ctx, f)
}

func (s *FeatureService) GetFeaturesByKeys(ctx context.Context, keys []string) ([]repository.FeatureRow, error) {
	if s.advRepo == nil {
		return nil, ErrNotImplemented
	}
	return s.advRepo.GetFeaturesByKeys(ctx, keys)
}

func (s *FeatureService) SearchFeaturesJSON(ctx context.Context, jsonExpr string, limit, offset int) ([]repository.FeatureRow, error) {
	if s.advRepo == nil {
		return nil, ErrNotImplemented
	}
	return s.advRepo.SearchFeaturesJSON(ctx, jsonExpr, limit, offset)
}
