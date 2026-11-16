package usecases

import (
	"context"
	"errors"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/domain"

	"golang.org/x/sync/singleflight"
)

type UrlRepository interface {
	CreateSlugIfNotExists(ctx context.Context, slugText string) (domain.Slug, error)
	CreateUrl(ctx context.Context, urlText string, slug domain.Slug) (domain.Url, error)
	GetUrlByText(ctx context.Context, urlText string) (domain.Url, error)
	GetUrlBySlug(ctx context.Context, slug domain.UrlSlug) (domain.Url, error)
}

type UrlCache interface {
	Save(ctx context.Context, url domain.Url) error
	GetUrlBySlug(ctx context.Context, slug domain.UrlSlug) (domain.Url, error)
	GetUrlByText(ctx context.Context, urlText string) (domain.Url, error)
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
func (uc *UrlUseCases) save(ctx context.Context, urlText string) (domain.Url, error) {
	// Проверка на случай если url уже существует
	_, err := uc.urlCache.GetUrlByText(ctx, urlText)
	if err == nil {
		return domain.Url{}, apperrors.ErrAlreadyExists
	}

	_, err = uc.urlRepository.GetUrlByText(ctx, urlText)
	if err == nil {
		return domain.Url{}, apperrors.ErrAlreadyExists
	}

	if errors.Is(err, apperrors.ErrUrlNotFound) {
		// Slug не найден, создаем новый и сохраняем в репозиторий

		slugText := generateSlug(6)
		slug, err := uc.urlRepository.CreateSlugIfNotExists(ctx, slugText)

		if err != nil {
			return domain.Url{}, err
		}

		url, err := uc.urlRepository.CreateUrl(ctx, urlText, slug)
		if err != nil {
			return domain.Url{}, err
		}

		return url, nil
	} else {
		// Неизвестная ошибка
		return domain.Url{}, err
	}
}

// Возвращает slug для переданного url.
//
// Если url не существует, то генерируется новый slug и сохраняется в бд
func (uc *UrlUseCases) GetUrlByText(ctx context.Context, urlText string) (domain.Url, error) {
	cachedUrl, err := uc.urlCache.GetUrlByText(ctx, urlText)
	if err == nil {
		return cachedUrl, nil
	}

	s, err, _ := uc.g.Do(urlText, func() (any, error) {
		url, err := uc.urlRepository.GetUrlByText(ctx, urlText)

		if err != nil {
			if errors.Is(err, apperrors.ErrUrlNotFound) {
				url, err = uc.save(ctx, urlText)
				if err != nil {
					return domain.Slug{}, err
				}
			} else {
				// Неизвестная ошибка
				return domain.Slug{}, err
			}
		}

		uc.urlCache.Save(ctx, url)
		return url, err
	})

	return s.(domain.Url), err
}

// Возвращает url для переданного slug
func (uc *UrlUseCases) GetUrlBySlug(ctx context.Context, slug domain.UrlSlug) (domain.Url, error) {
	url, err := uc.urlCache.GetUrlBySlug(ctx, slug)
	if err == nil {
		return url, nil
	}

	// Получаем Slug с repository, сохранение в cache, паттерн Single flight
	u, err, _ := uc.g.Do(slug.String(), func() (any, error) {
		url, err := uc.urlRepository.GetUrlBySlug(ctx, slug)
		if err != nil {
			return domain.Url{}, err
		}

		uc.urlCache.Save(ctx, url)
		return url, err
	})

	return u.(domain.Url), err
}
