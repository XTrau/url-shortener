package usecases

import (
	"context"
	"errors"
	"strconv"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/domain"

	"golang.org/x/sync/singleflight"
)

type UrlRepository interface {
	CreateSlugIfNotExists(ctx context.Context, slugText string) (domain.Slug, error)
	CreateUrl(ctx context.Context, urlText string, slug domain.Slug) (domain.Url, error)
	GetUrlBySlug(ctx context.Context, slug domain.UrlSlug) (domain.Url, error)
	GetSlugByUrl(ctx context.Context, urlText string) (domain.UrlSlug, error)
}

type UrlCache interface {
	Save(ctx context.Context, url string, slug domain.UrlSlug) error
	GetUrl(ctx context.Context, slug domain.UrlSlug) (string, error)
	GetSlug(ctx context.Context, url string) (domain.UrlSlug, error)
}

type UrlUseCases struct {
	urlRepository UrlRepository
	urlCache      UrlCache
	g             *singleflight.Group
}

func NewUrlUseCases(urlRepo UrlRepository, urlCache UrlCache) *UrlUseCases {
	return &UrlUseCases{
		urlRepository: urlRepo,
		urlCache:      urlCache,
		g:             &singleflight.Group{},
	}
}

// Сохраняет url в бд, создает для него новый slug
func (uc *UrlUseCases) save(ctx context.Context, urlText string) (domain.UrlSlug, error) {
	// Проверка на случай если url уже существует
	if _, err := uc.urlCache.GetSlug(ctx, urlText); err == nil {
		return domain.UrlSlug{}, apperrors.ErrAlreadyExists
	}

	_, err := uc.urlRepository.GetSlugByUrl(ctx, urlText)
	if err == nil {
		return domain.UrlSlug{}, apperrors.ErrAlreadyExists
	}

	if errors.Is(err, apperrors.ErrSlugNotFound) {
		// Slug не найден, создаем новый и сохраняем в репозиторий

		slugText := generateSlug(6)
		slug, err := uc.urlRepository.CreateSlugIfNotExists(ctx, slugText)

		if err != nil {
			return domain.UrlSlug{}, err
		}

		url, err := uc.urlRepository.CreateUrl(ctx, urlText, slug)
		if err != nil {
			return domain.UrlSlug{}, err
		}

		return *url.Slug, nil
	} else {
		// Неизвестная ошибка
		return domain.UrlSlug{}, err
	}
}

// Возвращает slug для переданного url.
// Если url не существует, то генерируется новый slug и сохраняется в бд
func (uc *UrlUseCases) GetSlug(ctx context.Context, urlText string) (domain.UrlSlug, error) {
	// Проверка в cache
	cachedSlug, err := uc.urlCache.GetSlug(ctx, urlText)
	if err == nil {
		return cachedSlug, nil
	}

	// Проверка в repository, сохранение в cache, паттерн Single flight
	s, err, _ := uc.g.Do(urlText, func() (any, error) {
		slug, err := uc.urlRepository.GetSlugByUrl(ctx, urlText)

		if err != nil {
			if errors.Is(err, apperrors.ErrSlugNotFound) {
				slug, err = uc.save(ctx, urlText)
				if err != nil {
					return domain.Slug{}, err
				}
			} else {
				// Неизвестная ошибка
				return domain.Slug{}, err
			}
		}

		uc.urlCache.Save(ctx, urlText, slug)
		return slug, err
	})

	return s.(domain.UrlSlug), err
}

// Возвращает url для переданного slug
func (uc *UrlUseCases) GetUrl(ctx context.Context, slug domain.UrlSlug) (domain.Url, error) {
	url, err := uc.urlCache.GetUrl(ctx, slug)
	if err == nil {
		return domain.Url{Url: url}, nil
	}

	// Group key для single flight
	groupKey := slug.Text + "_" + strconv.Itoa(slug.Number)

	// Получаем Slug с repository, сохранение в cache, паттерн Single flight
	u, err, _ := uc.g.Do(groupKey, func() (any, error) {
		url, err := uc.urlRepository.GetUrlBySlug(ctx, slug)
		if err != nil {
			return domain.Url{}, err
		}
		err = uc.urlCache.Save(ctx, url.Url, slug)
		return url, err
	})

	return u.(domain.Url), err
}
