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
	GetUrlByText(ctx context.Context, urlText string) (domain.Url, error)
	GetUrlBySlug(ctx context.Context, slug domain.UrlSlug) (domain.Url, error)
}

type RedirectUseCases interface {
	Save(ctx context.Context, redirect domain.RedirectInfo) error
}

type StatisticsUseCases interface {
	Summary(ctx context.Context, urlID int) (domain.SummaryStatistics, error)
	Timeline(ctx context.Context, urlID int, period domain.Period) (domain.TimelineStatistics, error)
}

type shortenerHandlers struct {
	urlUseCases        UrlUseCases
	redirectUseCases   RedirectUseCases
	statisticsUseCases StatisticsUseCases
}

func NewShortenerHandlers(urlUseCases UrlUseCases, redirectUseCases RedirectUseCases, statisticsUseCases StatisticsUseCases) *shortenerHandlers {
	return &shortenerHandlers{
		urlUseCases:        urlUseCases,
		redirectUseCases:   redirectUseCases,
		statisticsUseCases: statisticsUseCases,
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

	url, err := sr.urlUseCases.GetUrlByText(r.Context(), urlReq.Url)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(url.Slug); err != nil {
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

	url, err := sr.urlUseCases.GetUrlBySlug(r.Context(), urlSlug)

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

	if err = sr.redirectUseCases.Save(r.Context(), redirect); err != nil {
		slog.Error("Error on saving redirect info", slog.Any("redirect", redirect), slog.String("error", err.Error()))
	}

	w.Header().Set("Location", url.Url)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusMovedPermanently)
}

func (sr *shortenerHandlers) SummaryStatsHandler(w http.ResponseWriter, r *http.Request) {
	slugText := r.PathValue("slug")
	slog.Debug("Getting summary statistics.", slog.String("slug", slugText))

	urlSlug, err := domain.ParseUrlSlug(slugText)
	if err != nil {
		slog.Info("Error on parse url", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := sr.urlUseCases.GetUrlBySlug(r.Context(), urlSlug)
	if err != nil {
		if errors.Is(err, apperrors.ErrUrlNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			panic(err.Error())
		}
		return
	}

	stats, err := sr.statisticsUseCases.Summary(r.Context(), url.ID)
	stats.Url = url.Url

	if err != nil {
		slog.Info("Error getting stats for slug", slog.String("slug", slugText), slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&stats); err != nil {
		slog.Info("Error on encode stats", slog.String("slug", slugText), slog.Any("stats", stats), slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (sr *shortenerHandlers) PeriodStatsHandler(w http.ResponseWriter, r *http.Request) {
	slugText := r.PathValue("slug")
	slog.Debug("Getting summary statistics.", slog.String("slug", slugText))

	urlSlug, err := domain.ParseUrlSlug(slugText)
	if err != nil {
		slog.Info("Error on parse url", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	periodStr := r.URL.Query().Get("period")
	period, err := domain.ParsePeriod(periodStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := sr.urlUseCases.GetUrlBySlug(r.Context(), urlSlug)
	if err != nil {
		if errors.Is(err, apperrors.ErrUrlNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			panic(err.Error())
		}
		return
	}

	timelineStats, err := sr.statisticsUseCases.Timeline(r.Context(), url.ID, period)
	timelineStats.Url = url.Url

	if err != nil {
		slog.Info("Error getting stats for slug", slog.String("slug", slugText), slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&timelineStats); err != nil {
		slog.Info("Error on encode stats", slog.String("slug", slugText), slog.Any("stats", timelineStats), slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
