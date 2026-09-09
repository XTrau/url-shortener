package handlers

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ShortenerHandlers interface {
	ShortenerHandler(w http.ResponseWriter, r *http.Request)
	RedirectHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterRoutes(mux *http.ServeMux, sh ShortenerHandlers) {
	slog.Debug("Registering shortener routes")

	// Роуты для сокращения ссылок
	mux.HandleFunc("POST /short", sh.ShortenerHandler)
	mux.HandleFunc("GET /s/{slug}", sh.RedirectHandler)

	// Роуты для frontend
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	// Роут для Prometheus
	http.Handle("/metrics", promhttp.Handler())
}
