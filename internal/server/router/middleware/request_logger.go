package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func RequestLoggerMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		logger.Info(r.RequestURI,
			"method", r.Method,
			"duration", duration,
		)
	})
}
