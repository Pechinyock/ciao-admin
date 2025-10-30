package app

import (
	"ciao-admin/internal/loghandlers"
	"ciao-admin/internal/server/router/middleware"
	"log/slog"
	"net/http"
	"os"
)

func RegisterMiddlewares() (http.Handler, error) {
	slog.Info("starting configure middlewares")
	mux := http.NewServeMux()

	requestLogHandler := &loghandlers.PrettyStdLogHandler{
		Writer: os.Stdout,
		Level:  slog.LevelDebug,
	}
	requestLogger := slog.New(requestLogHandler)

	withLogging := middleware.RequestLoggerMiddleware(mux, requestLogger)

	return withLogging, nil
}
