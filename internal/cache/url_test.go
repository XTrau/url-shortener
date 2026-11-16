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

	url := domain.Url{
		ID:  1,
		Url: "http://example.com",
		Slug: domain.UrlSlug{
			Text:   "qweqwe",
			Number: 1,
		},
	}

	err := urlCache.Save(context.Background(), url)
	if err != nil {
		t.Errorf("urlCache.Save must return nil error")
	}

	cachedUrl, err := urlCache.GetUrlBySlug(context.Background(), url.Slug)

	if err != nil {
		t.Errorf("urlCache.GetByUrlSlug must return nil error")
	}

	if cachedUrl != url {
		t.Errorf("cachedURL=%q expected: %q", cachedUrl, url)
	}

	cachedUrl, err = urlCache.GetUrlByText(context.Background(), url.Url)

	if err != nil {
		t.Errorf("urlCache.GetByUrlSlug must return nil error")
	}

	if cachedUrl != url {
		t.Errorf("cachedURL=%q expected: %q", cachedUrl, url)
	}
}
