package repository

import (
	"context"
	"encoding/json"
	"errors"
)

var ErrNotFound = errors.New("feature not found")

type Repository interface {
	GetValueByKey(ctx context.Context, key string) (json.RawMessage, error)
	SetValueByKey(ctx context.Context, key string, value json.RawMessage) error
}
