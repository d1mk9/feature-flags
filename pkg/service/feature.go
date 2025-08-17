package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

// Vars — интерфейс, который используют хендлеры
type Vars interface {
	GetValue(key string) (json.RawMessage, error)
	SetValue(key string, value json.RawMessage) error
}

// внутренний элемент кэша с TTL
type cachedItem struct {
	val      json.RawMessage
	cachedAt time.Time
	ttl      time.Duration
}

// FeatureService — реализация Vars
type FeatureService struct {
	db    *sql.DB
	cache *lru.Cache[string, cachedItem]
	ttl   time.Duration
}

func NewFeatureService(db *sql.DB, cacheSize int, ttlMinutes int) *FeatureService {
	c, err := lru.New[string, cachedItem](cacheSize)
	if err != nil {
		panic(err)
	}
	return &FeatureService{
		db:    db,
		cache: c,
		ttl:   time.Duration(ttlMinutes) * time.Minute,
	}
}

func (s *FeatureService) GetValue(key string) (json.RawMessage, error) {
	if ci, ok := s.cache.Get(key); ok {
		if time.Since(ci.cachedAt) < ci.ttl {
			return ci.val, nil
		}
		s.cache.Remove(key) // TTL истёк
	}

	var value json.RawMessage
	err := s.db.QueryRow(`SELECT value FROM features WHERE key = $1`, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	s.cache.Add(key, cachedItem{val: value, cachedAt: time.Now(), ttl: s.ttl})
	return value, nil
}

func (s *FeatureService) SetValue(key string, value json.RawMessage) error {
	_, err := s.db.Exec(`
		INSERT INTO features (key, value, enabled)
		VALUES ($1, $2, true)
		ON CONFLICT (key) DO UPDATE
		SET value = EXCLUDED.value, updated_at = now()
	`, key, value)
	if err != nil {
		return err
	}
	s.cache.Remove(key) // инвалидация
	return nil
}
