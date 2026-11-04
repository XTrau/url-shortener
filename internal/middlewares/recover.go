package middlewares

import (
	"log/slog"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered.", slog.Any("error", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
