package cache_test

import (
	"context"
	"sync"
	"testing"
	"time"
	"urlshortener/internal/cache"
	"urlshortener/internal/domain"
)

// Struct without ttl for tests
type MapKeyValueStorage struct {
	cache map[string]string
	mu    sync.RWMutex
}

func (kv *MapKeyValueStorage) Get(ctx context.Context, key string) (string, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	value, ok := kv.cache[key]
	if !ok {
		return "", cache.ErrCacheKeyNotFound
	}
	return value, nil
}

func (kv *MapKeyValueStorage) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.cache[key] = value
	return nil
}

func TestUrlCache(t *testing.T) {
	cacheStore := &MapKeyValueStorage{
		cache: make(map[string]string),
		mu:    sync.RWMutex{},
	}
	urlCache := cache.NewUrlCache(cacheStore)

	url := "http://example.com"

	slug := domain.UrlSlug{
		Text:   "qweqwe",
		Number: 1,
	}

	err := urlCache.Save(context.TODO(), url, slug)
	if err != nil {
		t.Errorf("urlCache.Save must return nil error")
	}

	cachedSlug, err := urlCache.GetSlug(context.TODO(), url)

	if err != nil {
		t.Errorf("urlCache.GetSlug must return nil error")
	}

	if cachedSlug != slug {
		t.Errorf(
			"cachedSlug.Slug=%q, cachedSlug.SlugID=%d expected: cachedSlug.Slug=%q, cachedSlug.SlugID=%d",
			cachedSlug.Text,
			cachedSlug.Number,
			slug.Text,
			slug.Number,
		)
	}

	cachedURL, err := urlCache.GetUrl(context.TODO(), slug)

	if err != nil {
		t.Errorf("urlCache.GetUrl must return nil error")
	}

	if cachedURL != url {
		t.Errorf("cachedURL=%q expected: %q", cachedURL, url)
	}
}
