package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
	"urlshortener/internal/domain"
)

const UrlTTL = time.Minute * 5

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type UrlCache struct {
	cache Cache
}

func NewUrlCache(storage Cache) *UrlCache {
	return &UrlCache{storage}
}

func (uc *UrlCache) urlKey(slug string, slugID int) string {
	return fmt.Sprintf("url:%s#%d", slug, slugID)
}

func (uc *UrlCache) slugKey(url string) string {
	return fmt.Sprintf("slug:%s", url)
}

func (uc *UrlCache) Save(ctx context.Context, url string, slug domain.UrlSlug) error {
	urlKey := uc.urlKey(slug.Text, slug.Number)
	slugKey := uc.slugKey(url)

	slog.Debug(
		"Saving url to Cache",
		slog.String("urlKey", urlKey),
		slog.String("slugKey", slugKey),
	)

	slugData, err := json.Marshal(slug)

	if err != nil {
		return err
	}

	err = uc.cache.Set(ctx, slugKey, string(slugData), UrlTTL)
	if err != nil {
		return err
	}

	err = uc.cache.Set(ctx, urlKey, url, UrlTTL)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UrlCache) GetUrl(ctx context.Context, slug domain.UrlSlug) (string, error) {
	key := uc.urlKey(slug.Text, slug.Number)
	slog.Debug("Getting url from Cache", slog.String("key", key))

	return uc.cache.Get(ctx, key)
}

func (uc *UrlCache) GetSlug(ctx context.Context, url string) (domain.UrlSlug, error) {
	key := uc.slugKey(url)
	slog.Debug("Getting slug from Redis", slog.String("key", key))

	slugData, err := uc.cache.Get(ctx, key)

	if err != nil {
		return domain.UrlSlug{}, err
	}

	var cachedSlug domain.UrlSlug
	if err := json.Unmarshal([]byte(slugData), &cachedSlug); err != nil {
		return domain.UrlSlug{}, err
	}

	return cachedSlug, nil
}
