package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/domain"
)

type UrlUseCases interface {
	GetSlug(ctx context.Context, urlText string) (domain.UrlSlug, error)
	GetUrl(ctx context.Context, slug domain.UrlSlug) (domain.Url, error)
}

type shortenerHandlers struct {
	useCases UrlUseCases
}

func NewShortenerHandlers(useCases UrlUseCases) *shortenerHandlers {
	return &shortenerHandlers{useCases}
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

	slug, err := sr.useCases.GetSlug(r.Context(), urlReq.Url)
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

	url, err := sr.useCases.GetUrl(r.Context(), urlSlug)

	if err != nil {
		if errors.Is(err, apperrors.ErrUrlNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			panic(err)
		}
		return
	}

	w.Header().Set("Location", url.Url)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusMovedPermanently)
}
