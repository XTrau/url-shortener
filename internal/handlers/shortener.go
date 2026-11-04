package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/domain"
)

type UrlUseCases interface {
	GetSlug(url domain.Url) (domain.Slug, error)
	GetUrl(slug domain.Slug) (domain.Url, error)
}

type shortenerHandlers struct {
	useCases UrlUseCases
}

func NewShortenerHandlers(useCases UrlUseCases) *shortenerHandlers {
	return &shortenerHandlers{useCases}
}

func (sr *shortenerHandlers) ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	var urlReq domain.Url

	if err := json.NewDecoder(r.Body).Decode(&urlReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	slug, err := sr.useCases.GetSlug(urlReq)
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
	slugID, err := strconv.Atoi(r.URL.Query().Get("i"))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	slug := domain.Slug{
		SlugID: slugID,
		Text:   slugText,
	}

	url, err := sr.useCases.GetUrl(slug)

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
