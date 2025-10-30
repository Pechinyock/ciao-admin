package app

import (
	"ciao-admin/internal/loghandlers"
	"ciao-admin/internal/server"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type AdminApplication struct {
	Version       string
	configuration server.CoreConfig
}

func (app *AdminApplication) Init() bool {
	handler := loghandlers.PrettyStdLogHandler{
		Level:  slog.LevelDebug,
		Writer: os.Stdout,
	}

	logger := slog.New(&handler)
	slog.SetDefault(logger)
	slog.Info("start initializing admin application", "version", app.Version)

	configPath := readFlags()
	if configPath == DefaultConfigPath {
		slog.Warn("config file was not provided, trying to load the default one",
			"default config file paht", DefaultConfigPath,
		)
	}

	cfg := loadConfigFromFile(configPath)
	if cfg == nil {
		return false
	}

	if !verifyServerConfig(cfg) {
		return false
	}

	app.configuration = *cfg
	configuredLogHandler := configureLogging(cfg.LogLevel, cfg.LogDirPath)
	completeLogger := slog.New(configuredLogHandler)
	slog.SetDefault(completeLogger)
	return true
}

func (app *AdminApplication) Run() {
	fullAddr := fmt.Sprintf("%s:%d", app.configuration.Address,
		app.configuration.Port,
	)

	mux, err := RegisterMiddlewares()

	if err != nil {
		slog.Error("failed to configure middlewares")
		return
	}
	serv := http.Server{
		Addr:         fullAddr,
		ReadTimeout:  app.configuration.ReadTimeout * time.Second,
		WriteTimeout: app.configuration.WriteTimeout * time.Second,
		Handler:      mux,
	}

	isHttps := app.configuration.CertFilePath != ""

	if isHttps {
		slog.Info(fmt.Sprintf("running website at https://%s", fullAddr))
		err := serv.ListenAndServeTLS(app.configuration.CertFilePath,
			app.configuration.CertKeyFilePath,
		)
		if err != nil {
			slog.Error("failed to run server", "reason", err.Error())
		}
	} else {
		slog.Info(fmt.Sprintf("running website at http://%s", fullAddr))
		err := serv.ListenAndServe()
		if err != nil {
			slog.Error("failed to run server", "reason", err.Error())
		}
	}
}
