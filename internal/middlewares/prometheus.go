package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество HTTP запросов.",
		},
		[]string{"method", "endpoint", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Время ответа на HTTP запрос в секундах.",
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (srw *responseWriter) WriteHeader(statusCode int) {
	srw.statusCode = statusCode
	srw.ResponseWriter.WriteHeader(statusCode)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{w, 200}

		start := time.Now()
		next.ServeHTTP(rw, r)

		// Время исполнения запроса
		elapsed := time.Since(start)

		// Записываем время
		requestDuration.WithLabelValues(
			r.Method,
			r.Pattern,
		).Observe(elapsed.Seconds())

		// Увеличиваем счетчик запросов
		requestsTotal.WithLabelValues(
			r.Method,
			r.Pattern,
			strconv.Itoa(rw.statusCode),
		).Inc()
	})
}
