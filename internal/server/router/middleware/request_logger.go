package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

type RequestLogger struct {
	Logger slog.Logger
}

func (rl *RequestLogger) Add(next http.Handler) http.Handler {
	slog.Info("request logger middleware has been added")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		if rw.statusCode >= http.StatusOK && rw.statusCode < http.StatusMultipleChoices {
			rl.Logger.Info(r.RequestURI,
				"method", r.Method,
				"status code", rw.statusCode,
				"duration", duration,
			)
		} else if rw.statusCode >= http.StatusBadRequest && rw.statusCode <= http.StatusNetworkAuthenticationRequired {
			rl.Logger.Error(r.RequestURI,
				"method", r.Method,
				"status code", rw.statusCode,
				"duration", duration,
			)
		}
	})
}
