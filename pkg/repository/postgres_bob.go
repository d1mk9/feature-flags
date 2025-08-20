package repository

import (
	"context"
)

func (r *PostgresRepository) ListFeatures(ctx context.Context, f FeatureListFilters) ([]FeatureRow, error) {
	return nil, ErrNotImplemented
}

func (r *PostgresRepository) CountFeatures(ctx context.Context, f FeatureListFilters) (int64, error) {
	return 0, ErrNotImplemented
}

func (r *PostgresRepository) GetFeaturesByKeys(ctx context.Context, keys []string) ([]FeatureRow, error) {
	return nil, ErrNotImplemented
}

func (r *PostgresRepository) SearchFeaturesJSON(ctx context.Context, jsonExpr string, limit, offset int) ([]FeatureRow, error) {
	return nil, ErrNotImplemented
}
