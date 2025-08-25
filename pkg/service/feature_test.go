package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"feature-flags/pkg/repository"
	"feature-flags/pkg/service"
)

type fakeRepo struct {
	vals map[string]json.RawMessage
	err  error
}

func (f *fakeRepo) GetValueByKey(ctx context.Context, key string) (json.RawMessage, error) {
	if f.err != nil {
		return nil, f.err
	}
	v, ok := f.vals[key]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return v, nil
}
func (f *fakeRepo) SetValueByKey(ctx context.Context, key string, value json.RawMessage) error {
	if f.err != nil {
		return f.err
	}
	f.vals[key] = value
	return nil
}

func TestFeatureService_TTLExpires(t *testing.T) {
	repo := &fakeRepo{vals: map[string]json.RawMessage{"k": json.RawMessage(`1`)}}
	s, _ := service.NewFeatureService(repo, 16, 0)
	ctx := context.Background()

	_, _ = s.GetValue(ctx, "k")
	repo.err = errors.New("after first fetch")
	_, err := s.GetValue(ctx, "k")
	if err == nil {
		t.Fatalf("expected error after TTL expire")
	}
}
