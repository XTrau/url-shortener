package usecases

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/domain"

	"golang.org/x/sync/singleflight"
)

const charSet = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM134567890"

func generateSlug(size int) string {
	var b strings.Builder
	b.Grow(size)

	for range size {
		ind := rand.Intn(len(charSet))
		b.Write([]byte(charSet[ind : ind+1]))
	}
	return b.String()
}

type UrlRepository interface {
	Create(url domain.Url, slug domain.Slug) error
	GetUrlBySlug(slug domain.Slug) (domain.Url, error)
	GetSlugByUrl(url domain.Url) (domain.Slug, error)
	GetFreeSlugID(slugText string) (int, error)
}

type UrlCache interface {
	Save(url domain.Url, slug domain.Slug) error
	GetUrl(slug domain.Slug) (domain.Url, error)
	GetSlug(url domain.Url) (domain.Slug, error)
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
func (uc *UrlUseCases) save(url domain.Url) (domain.Slug, error) {
	var slug domain.Slug

	// Проверка на случай если url уже существует
	if _, err := uc.urlCache.GetSlug(url); err == nil {
		return domain.Slug{}, apperrors.ErrAlreadyExists
	}

	_, err := uc.urlRepository.GetSlugByUrl(url)
	if err == nil {
		return domain.Slug{}, apperrors.ErrAlreadyExists
	}

	if errors.Is(err, apperrors.ErrSlugNotFound) {
		// Slug не найден, создаем новый и сохраняем в репозиторий

		slugText := generateSlug(8)
		slugID, err := uc.urlRepository.GetFreeSlugID(slugText)

		if err != nil {
			return domain.Slug{}, err
		}

		slug = domain.Slug{
			SlugID: slugID,
			Text:   slugText,
		}

		err = uc.urlRepository.Create(url, slug)
		if err != nil {
			return domain.Slug{}, err
		}

		return slug, nil
	} else {
		// Неизвестная ошибка
		return domain.Slug{}, err
	}
}

// Возвращает slug для переданного url.
// Если url не существует, то генерируется новый slug и сохраняется в бд
func (uc *UrlUseCases) GetSlug(url domain.Url) (domain.Slug, error) {
	// Проверка в cache
	cachedSlug, err := uc.urlCache.GetSlug(url)
	if err == nil {
		return cachedSlug, nil
	}

	// Проверка в repository, сохранение в cache, паттерн Single flight
	s, err, _ := uc.g.Do(url.Url, func() (any, error) {
		slug, err := uc.urlRepository.GetSlugByUrl(url)

		if err != nil {
			if errors.Is(err, apperrors.ErrSlugNotFound) {
				slug, err = uc.save(url)
				if err != nil {
					return domain.Slug{}, err
				}
			} else {
				// Неизвестная ошибка
				return domain.Slug{}, err
			}
		}

		uc.urlCache.Save(url, slug)
		return slug, err
	})

	return s.(domain.Slug), err
}

// Возвращает url для переданного slug
func (uc *UrlUseCases) GetUrl(slug domain.Slug) (domain.Url, error) {
	url, err := uc.urlCache.GetUrl(slug)
	if err == nil {
		return url, nil
	}

	// Group key для single flight
	groupKey := slug.Text + "_" + strconv.Itoa(slug.SlugID)

	// Получаем Slug с repository, сохранение в cache, паттерн Single flight
	u, err, _ := uc.g.Do(groupKey, func() (any, error) {
		url, err := uc.urlRepository.GetUrlBySlug(slug)
		if err != nil {
			return domain.Url{}, err
		}
		err = uc.urlCache.Save(url, slug)
		return url, err
	})

	return u.(domain.Url), err
}
