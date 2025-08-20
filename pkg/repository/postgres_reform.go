package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"feature-flags/pkg/models"

	reform "gopkg.in/reform.v1"
)

type PostgresRepository struct {
	db *reform.DB
}

func NewPostgresRepository(rdb *reform.DB) *PostgresRepository {
	return &PostgresRepository{db: rdb}
}

func (r *PostgresRepository) GetValue(ctx context.Context, key string) (json.RawMessage, error) {
	rec, err := r.db.WithContext(ctx).FindOneFrom(models.FeaturesTable, "key", key)
	if err != nil {
		if errors.Is(err, reform.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo.GetValue find key=%q: %w", key, err)
	}
	return rec.(*models.Features).Value, nil
}

func (r *PostgresRepository) UpsertValue(ctx context.Context, key string, value json.RawMessage) error {
	rec, err := r.db.WithContext(ctx).FindOneFrom(models.FeaturesTable, "key", key)
	switch {
	case err == nil:
		f := rec.(*models.Features)
		f.Value = value
		if err := r.db.WithContext(ctx).Update(f); err != nil {
			return fmt.Errorf("repo.UpsertValue update key=%q: %w", key, err)
		}
		return nil

	case errors.Is(err, reform.ErrNoRows):
		if err := r.db.WithContext(ctx).Insert(&models.Features{
			Key:     key,
			Enabled: true,
			Value:   value,
		}); err != nil {
			return fmt.Errorf("repo.UpsertValue insert key=%q: %w", key, err)
		}
		return nil

	default:
		return fmt.Errorf("repo.UpsertValue select key=%q: %w", key, err)
	}
}
