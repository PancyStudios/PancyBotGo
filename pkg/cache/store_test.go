package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubStore struct {
	get func(context.Context, string) ([]byte, error)
}

func (s stubStore) Get(ctx context.Context, key string) ([]byte, error) {
	return s.get(ctx, key)
}
func (stubStore) Set(context.Context, string, []byte, time.Duration) error { return nil }
func (stubStore) Delete(context.Context, string) error                     { return nil }
func (stubStore) Clear(context.Context) error                              { return nil }
func (stubStore) Size(context.Context) (int, error)                        { return 0, nil }
func (stubStore) Close() error                                             { return nil }
func (stubStore) Name() string                                             { return "stub" }

func TestMemoryStoreCopiesValues(t *testing.T) {
	store := NewMemoryStore(10)
	input := []byte("original")
	if err := store.Set(context.Background(), "key", input, time.Minute); err != nil {
		t.Fatal(err)
	}
	input[0] = 'X'

	value, err := store.Get(context.Background(), "key")
	if err != nil {
		t.Fatal(err)
	}
	value[0] = 'Y'
	again, err := store.Get(context.Background(), "key")
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != "original" {
		t.Fatalf("cached bytes were mutated: %q", again)
	}
}

func TestMemoryStoreExpiresAndEvicts(t *testing.T) {
	store := NewMemoryStore(1)
	if err := store.Set(context.Background(), "old", []byte("1"), time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := store.Set(context.Background(), "new", []byte("2"), time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(context.Background(), "old"); !errors.Is(err, ErrMiss) {
		t.Fatalf("expected evicted key to miss, got %v", err)
	}

	if err := store.Set(context.Background(), "short", []byte("3"), time.Millisecond); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond)
	if _, err := store.Get(context.Background(), "short"); !errors.Is(err, ErrMiss) {
		t.Fatalf("expected expired key to miss, got %v", err)
	}
}

func TestMemoryStoreClear(t *testing.T) {
	store := NewMemoryStore(10)
	_ = store.Set(context.Background(), "key", []byte("value"), 0)
	if err := store.Clear(context.Background()); err != nil {
		t.Fatal(err)
	}
	if size, err := store.Size(context.Background()); err != nil || size != 0 {
		t.Fatalf("size after clear = %d, err = %v", size, err)
	}
}

func TestFailoverStoreUsesLocalCacheOnlyOnPrimaryError(t *testing.T) {
	ctx := context.Background()
	local := NewMemoryStore(10)
	_ = local.Set(ctx, "key", []byte("fallback"), time.Minute)

	down := stubStore{get: func(context.Context, string) ([]byte, error) {
		return nil, errors.New("redis unavailable")
	}}
	store := NewFailoverStore(down, local)
	value, err := store.Get(ctx, "key")
	if err != nil || string(value) != "fallback" {
		t.Fatalf("fallback value = %q, err = %v", value, err)
	}

	missing := stubStore{get: func(context.Context, string) ([]byte, error) {
		return nil, ErrMiss
	}}
	store = NewFailoverStore(missing, local)
	if _, err := store.Get(ctx, "key"); !errors.Is(err, ErrMiss) {
		t.Fatalf("normal primary miss must not serve stale local data: %v", err)
	}
	if _, err := local.Get(ctx, "key"); !errors.Is(err, ErrMiss) {
		t.Fatalf("normal primary miss did not invalidate local fallback: %v", err)
	}
}

func TestRedisStoreRejectsUnsafePrefixBeforeConnecting(t *testing.T) {
	if _, err := NewRedisStore(context.Background(), "redis://localhost:6379/0", ""); err == nil {
		t.Fatal("empty Redis prefix was accepted")
	}
	if _, err := NewRedisStore(context.Background(), "redis://localhost:6379/0", "pancybot:*"); err == nil {
		t.Fatal("wildcard Redis prefix was accepted")
	}
}
