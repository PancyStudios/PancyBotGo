// Package cache provides interchangeable cache backends for shared bot data.
package cache

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

// ErrMiss indicates that a key does not exist or has expired.
var ErrMiss = errors.New("cache miss")

// Store is the minimal cache contract used by the database layer.
// Cached values are disposable; MongoDB remains the source of truth.
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	Size(ctx context.Context) (int, error)
	Close() error
	Name() string
}

type memoryEntry struct {
	key       string
	value     []byte
	expiresAt time.Time
}

// MemoryStore is a bounded, process-local LRU cache.
type MemoryStore struct {
	mu         sync.Mutex
	maxEntries int
	entries    map[string]*list.Element
	lru        *list.List
}

// NewMemoryStore creates an in-memory cache. A non-positive limit disables eviction.
func NewMemoryStore(maxEntries int) *MemoryStore {
	return &MemoryStore{
		maxEntries: maxEntries,
		entries:    make(map[string]*list.Element),
		lru:        list.New(),
	}
}

func cloneBytes(value []byte) []byte {
	return append([]byte(nil), value...)
}

func (m *MemoryStore) Get(ctx context.Context, key string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	element, ok := m.entries[key]
	if !ok {
		return nil, ErrMiss
	}
	entry := element.Value.(*memoryEntry)
	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		m.removeElement(element)
		return nil, ErrMiss
	}

	m.lru.MoveToFront(element)
	return cloneBytes(entry.value), nil
}

func (m *MemoryStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if element, ok := m.entries[key]; ok {
		entry := element.Value.(*memoryEntry)
		entry.value = cloneBytes(value)
		entry.expiresAt = expiresAt
		m.lru.MoveToFront(element)
		return nil
	}

	element := m.lru.PushFront(&memoryEntry{key: key, value: cloneBytes(value), expiresAt: expiresAt})
	m.entries[key] = element
	if m.maxEntries > 0 && m.lru.Len() > m.maxEntries {
		m.removeElement(m.lru.Back())
	}
	return nil
}

func (m *MemoryStore) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if element, ok := m.entries[key]; ok {
		m.removeElement(element)
	}
	return nil
}

func (m *MemoryStore) Clear(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = make(map[string]*list.Element)
	m.lru.Init()
	return nil
}

func (m *MemoryStore) Size(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.removeExpired(time.Now())
	return len(m.entries), nil
}

func (m *MemoryStore) Close() error { return nil }
func (m *MemoryStore) Name() string { return "memory" }

func (m *MemoryStore) removeElement(element *list.Element) {
	if element == nil {
		return
	}
	entry := element.Value.(*memoryEntry)
	delete(m.entries, entry.key)
	m.lru.Remove(element)
}

func (m *MemoryStore) removeExpired(now time.Time) {
	for element := m.lru.Back(); element != nil; {
		previous := element.Prev()
		entry := element.Value.(*memoryEntry)
		if !entry.expiresAt.IsZero() && now.After(entry.expiresAt) {
			m.removeElement(element)
		}
		element = previous
	}
}

// FailoverStore uses the primary cache normally and a local cache only when the
// primary is unavailable. It never serves the local copy after a normal primary miss.
type FailoverStore struct {
	primary  Store
	fallback Store
}

func NewFailoverStore(primary, fallback Store) *FailoverStore {
	return &FailoverStore{primary: primary, fallback: fallback}
}

func (f *FailoverStore) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := f.primary.Get(ctx, key)
	if err == nil {
		return value, nil
	}
	if errors.Is(err, ErrMiss) {
		_ = f.fallback.Delete(context.Background(), key)
		return nil, ErrMiss
	}
	return f.fallback.Get(ctx, key)
}

func (f *FailoverStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	_ = f.fallback.Set(ctx, key, value, ttl)
	return f.primary.Set(ctx, key, value, ttl)
}

func (f *FailoverStore) Delete(ctx context.Context, key string) error {
	_ = f.fallback.Delete(ctx, key)
	return f.primary.Delete(ctx, key)
}

func (f *FailoverStore) Clear(ctx context.Context) error {
	_ = f.fallback.Clear(ctx)
	return f.primary.Clear(ctx)
}

func (f *FailoverStore) Size(ctx context.Context) (int, error) {
	return f.primary.Size(ctx)
}

func (f *FailoverStore) Close() error {
	primaryErr := f.primary.Close()
	fallbackErr := f.fallback.Close()
	return errors.Join(primaryErr, fallbackErr)
}

func (f *FailoverStore) Name() string { return f.primary.Name() + "+memory-fallback" }
