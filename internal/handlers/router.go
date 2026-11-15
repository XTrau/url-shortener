package handlers

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ShortenerHandlers interface {
	ShortenerHandler(w http.ResponseWriter, r *http.Request)
	RedirectHandler(w http.ResponseWriter, r *http.Request)
	SummaryStatsHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterRoutes(mux *http.ServeMux, sh ShortenerHandlers) {
	slog.Debug("Registering shortener routes")

	// Роуты для сокращения ссылок
	mux.HandleFunc("POST /api/short", sh.ShortenerHandler)
	mux.HandleFunc("GET /s/{slug}", sh.RedirectHandler)

	// Роуты для статистики
	mux.HandleFunc("GET /api/statistics", sh.SummaryStatsHandler)

	// Роуты для frontend
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	// Роут для Prometheus
	mux.Handle("/metrics", promhttp.Handler())
}
