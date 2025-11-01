package app

import (
	"ciao-admin/internal/loghandlers"
	"ciao-admin/internal/server"
	"ciao-admin/internal/server/router/middleware"
	"log/slog"
	"net/http"
	"os"
)

func RegisterMiddlewares(config *server.MeddlewaresConfig) (http.Handler, error) {
	slog.Info("starting configure middlewares")
	mux := http.NewServeMux()

	if config == nil {
		slog.Warn("middleware configuration is not provided, enable all of them...")
		middlewares := getAll()
		result := registerMiddlewares(middlewares, mux)
		return result, nil
	}

	var middlewares []middleware.Middleware
	rqLoggerConfig := config.RequestLogger
	middlewares = append(middlewares, getRequestLogger(rqLoggerConfig))

	var result = registerMiddlewares(middlewares, mux)
	/*[STOPS HERE]*/
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return result, nil
}

func getAll() []middleware.Middleware {
	var middleares []middleware.Middleware
	middleares = append(middleares, getRequestLogger(nil))
	return middleares
}

func registerMiddlewares(middlewares []middleware.Middleware, mux *http.ServeMux) http.Handler {
	var result http.Handler
	for _, mw := range middlewares {
		if mw == nil {
			continue
		}
		result = mw.Add(mux)
	}
	return result
}

func getRequestLogger(config *server.RequestLoggerMiddlewareConfig) middleware.Middleware {
	var handler slog.Handler
	if config == nil {
		handler = &loghandlers.PrettyStdLogHandler{
			Writer: os.Stdout,
			Level:  slog.LevelInfo,
		}
	} else {
		if !config.Enabled {
			return nil
		}
		logLevel := parseLogLevel(config.LogLevel)
		handler = &loghandlers.PrettyStdLogHandler{
			Writer: os.Stdout,
			Level:  logLevel,
		}
	}
	requestLogger := slog.New(handler)
	rqLogger := middleware.RequestLogger{
		Logger: *requestLogger,
	}
	return &rqLogger
}
