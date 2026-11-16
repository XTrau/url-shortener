package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"urlshortener/internal/apperrors"
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

func (uc *UrlCache) Save(ctx context.Context, url domain.Url) error {
	urlKey := uc.slugKey(url.Slug.Text, url.Slug.Number)
	slugKey := uc.urlKey(url.Url)

	slog.Debug(
		"Saving url to Cache",
		slog.String("urlKey", urlKey),
		slog.String("slugKey", slugKey),
	)

	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	err = uc.cache.Set(ctx, slugKey, string(data), UrlTTL)
	if err != nil {
		return err
	}

	err = uc.cache.Set(ctx, urlKey, string(data), UrlTTL)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UrlCache) GetUrlBySlug(ctx context.Context, slug domain.UrlSlug) (domain.Url, error) {
	key := uc.slugKey(slug.Text, slug.Number)
	slog.Debug("Getting url from Cache", slog.String("key", key))

	urlData, err := uc.cache.Get(ctx, key)
	if err != nil {
		switch {
		case errors.Is(err, ErrCacheKeyNotFound):
			return domain.Url{}, apperrors.ErrUrlNotFound
		default:
			return domain.Url{}, err
		}
	}

	var url domain.Url
	if err := json.Unmarshal([]byte(urlData), &url); err != nil {
		return domain.Url{}, err
	}

	return url, nil
}

func (uc *UrlCache) GetUrlByText(ctx context.Context, urlText string) (domain.Url, error) {
	key := uc.urlKey(urlText)
	slog.Debug("Getting url from cache", slog.String("key", key))

	urlData, err := uc.cache.Get(ctx, key)
	if err != nil {
		switch {
		case errors.Is(err, ErrCacheKeyNotFound):
			return domain.Url{}, apperrors.ErrUrlNotFound
		default:
			return domain.Url{}, err
		}
	}

	var url domain.Url
	if err := json.Unmarshal([]byte(urlData), &url); err != nil {
		return domain.Url{}, err
	}

	return url, nil
}

func (uc *UrlCache) slugKey(slug string, slugID int) string {
	return fmt.Sprintf("url:%s#%d", slug, slugID)
}

func (uc *UrlCache) urlKey(url string) string {
	return fmt.Sprintf("slug:%s", url)
}
