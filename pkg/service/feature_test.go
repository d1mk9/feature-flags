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

func (f *fakeRepo) GetValue(ctx context.Context, key string) (json.RawMessage, error) {
	if f.err != nil {
		return nil, f.err
	}
	v, ok := f.vals[key]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return v, nil
}
func (f *fakeRepo) UpsertValue(ctx context.Context, key string, value json.RawMessage) error {
	if f.err != nil {
		return f.err
	}
	f.vals[key] = value
	return nil
}

func TestFeatureService_CacheHit(t *testing.T) {
	repo := &fakeRepo{vals: map[string]json.RawMessage{"k": json.RawMessage(`123`)}}
	s, err := service.NewFeatureService(repo, 16, 15)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	got1, err := s.GetValue(ctx, "k")
	if err != nil {
		t.Fatalf("GetValue: %v", err)
	}
	repo.err = errors.New("boom")
	got2, err := s.GetValue(ctx, "k")
	if err != nil {
		t.Fatalf("GetValue (cached): %v", err)
	}
	if string(got1) != string(got2) {
		t.Fatalf("cache miss: %s vs %s", got1, got2)
	}
}

func TestFeatureService_SetValue_InvalidatesCache(t *testing.T) {
	repo := &fakeRepo{vals: map[string]json.RawMessage{"k": json.RawMessage(`1`)}}
	s, _ := service.NewFeatureService(repo, 16, 15)
	ctx := context.Background()

	_, _ = s.GetValue(ctx, "k")
	if err := s.SetValue(ctx, "k", json.RawMessage(`2`)); err != nil {
		t.Fatalf("SetValue: %v", err)
	}

	got, err := s.GetValue(ctx, "k")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "2" {
		t.Fatalf("want 2, got %s", got)
	}
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
