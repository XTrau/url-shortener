package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/domain"
)

type UrlUseCases interface {
	GetSlug(ctx context.Context, urlText string) (domain.UrlSlug, error)
	GetUrl(ctx context.Context, slug domain.UrlSlug) (domain.Url, error)
}

type RedirectUseCases interface {
	Save(ctx context.Context, redirect domain.RedirectInfo) error
}

type shortenerHandlers struct {
	urlUseCases      UrlUseCases
	redirectUseCases RedirectUseCases
}

func NewShortenerHandlers(urlUseCases UrlUseCases, redirectUseCases RedirectUseCases) *shortenerHandlers {
	return &shortenerHandlers{
		urlUseCases:      urlUseCases,
		redirectUseCases: redirectUseCases,
	}
}

type GetSlugRequest struct {
	Url string `json:"url"`
}

func (sr *shortenerHandlers) ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	var urlReq GetSlugRequest

	if err := json.NewDecoder(r.Body).Decode(&urlReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	slug, err := sr.urlUseCases.GetSlug(r.Context(), urlReq.Url)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slug); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		panic(err)
	}
}

func (sr *shortenerHandlers) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	slugText := r.PathValue("slug")
	slog.Debug("Getting redirect on slug", slog.String("slug", slugText))

	urlSlug, err := domain.ParseUrlSlug(slugText)

	if err != nil {
		slog.Info("Error on parse url", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := sr.urlUseCases.GetUrl(r.Context(), urlSlug)

	if err != nil {
		if errors.Is(err, apperrors.ErrUrlNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			panic(err)
		}
		return
	}

	redirect := domain.RedirectInfo{
		CreatedAt: time.Now(),
		IP:        r.RemoteAddr,
		UserAgent: r.UserAgent(),
		UrlID:     url.ID,
	}

	sr.redirectUseCases.Save(r.Context(), redirect)

	w.Header().Set("Location", url.Url)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusMovedPermanently)
}

// func (sr *shortenerHandlers) SummaryStatsHandler(w http.ResponseWriter, r *http.Request) {
// 	slugText := r.PathValue("slug")
// 	slog.Debug("Getting summary statistics.", slog.String("slug", slugText))

// 	urlSlug, err := domain.ParseUrlSlug(slugText)
// 	if err != nil {
// 		slog.Info("Error on parse url", slog.String("error", err.Error()))
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	stats, err := sr.useCases.Statistics(r.Context(), urlSlug)

// 	if err != nil {
// 		slog.Info("Error getting stats for slug", slog.String("slug", slugText), slog.String("error", err.Error()))
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}

// 	if err := json.NewEncoder(w).Encode(&stats); err != nil {
// 		slog.Info("Error on encode stats", slog.String("slug", slugText), slog.Any("stats", stats), slog.String("error", err.Error()))
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}
// }
