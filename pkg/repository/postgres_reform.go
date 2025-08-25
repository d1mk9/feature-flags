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

func (r *PostgresRepository) GetValueByKey(ctx context.Context, key string) (json.RawMessage, error) {
	var f models.Features
	if err := r.db.WithContext(ctx).FindOneTo(&f, "key", key); err != nil {
		if errors.Is(err, reform.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo.GetValue find key=%q: %w", key, err)
	}
	return f.Value, nil
}

func (r *PostgresRepository) SetValueByKey(ctx context.Context, key string, value json.RawMessage) error {
	var f models.Features
	err := r.db.WithContext(ctx).FindOneTo(&f, "key", key)
	switch {
	case err == nil:
		f.Value = value
		if err := r.db.WithContext(ctx).Update(&f); err != nil {
			return fmt.Errorf("repo.UpsertValue update key=%q: %w", key, err)
		}
		return nil

	case errors.Is(err, reform.ErrNoRows):
		f = models.Features{
			Key:     key,
			Enabled: true,
			Value:   value,
		}
		if err := r.db.WithContext(ctx).Insert(&f); err != nil {
			return fmt.Errorf("repo.UpsertValue insert key=%q: %w", key, err)
		}
		return nil

	default:
		return fmt.Errorf("repo.UpsertValue select key=%q: %w", key, err)
	}
}
