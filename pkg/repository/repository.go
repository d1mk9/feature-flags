package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type Repository interface {
	GetValue(ctx context.Context, key string) (json.RawMessage, error)
	UpsertValue(ctx context.Context, key string, value json.RawMessage) error
}

var (
	ErrNotFound       = errors.New("feature not found")
	ErrNotImplemented = errors.New("not implemented")
)

type FeatureRow struct {
	ID          string
	Key         string
	Description *string
	Enabled     bool
	Value       json.RawMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FeatureListFilters struct {
	Enabled      *bool
	KeyLike      string
	JSONContains string
	Limit        int
	Offset       int
	OrderBy      string
}

type AdvancedRepository interface {
	ListFeatures(ctx context.Context, f FeatureListFilters) ([]FeatureRow, error)
	CountFeatures(ctx context.Context, f FeatureListFilters) (int64, error)
	GetFeaturesByKeys(ctx context.Context, keys []string) ([]FeatureRow, error)
	SearchFeaturesJSON(ctx context.Context, jsonExpr string, limit, offset int) ([]FeatureRow, error)
}
