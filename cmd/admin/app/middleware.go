package app

import (
	"ciao-admin/cmd/admin/app/configuration"
	"ciao-admin/internal/loghandlers"
	"ciao-admin/internal/server/router/middleware"
	"errors"
	"log/slog"
	"net/http"
	"os"
)

func registerMiddlewares(router *http.ServeMux, config *configuration.MeddlewaresConfig) (http.Handler, error) {
	if router == nil {
		return nil, errors.New("trying to add middlwares to nil router")
	}
	slog.Info("starting configure middlewares")

	mux := router

	middlewares := getNonDisablable()

	if config == nil {
		slog.Warn("middleware configuration is not provided, enable all of them...")
		allMiddleWares := append(middlewares, getAll()...)
		result := addMiddlewares(allMiddleWares, mux)
		return result, nil
	}

	rqLoggerConfig := config.RequestLogger
	middlewares = append(middlewares, getRequestLogger(rqLoggerConfig))

	var result = addMiddlewares(middlewares, mux)

	return result, nil
}

func getAll() []middleware.Middleware {
	var middleares []middleware.Middleware
	middleares = append(middleares, getRequestLogger(nil))
	return middleares
}

func addMiddlewares(middlewares []middleware.Middleware, mux *http.ServeMux) http.Handler {
	var result http.Handler
	for _, mw := range middlewares {
		if mw == nil {
			continue
		}
		result = mw.Add(mux)
	}
	return result
}

func getRequestLogger(config *configuration.RequestLoggerMiddlewareConfig) middleware.Middleware {
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

func getNonDisablable() []middleware.Middleware {
	result := []middleware.Middleware{
		&middleware.ServerRecovery{},
	}
	return result
}
