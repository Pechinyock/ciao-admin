package app

import (
	"ciao-admin/cmd/coven/app/configuration"
	"ciao-admin/internal/UI/webui"
	"ciao-admin/internal/loghandlers"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type CovenApplication struct {
	Version          string
	configuration    *configuration.CoreConfig
	middlewareConfig *configuration.MeddlewaresConfig
	sharFileConfig   *configuration.FileServerConfig
	webUIConfig      *webui.WebUIBundleConfig
	webUIBundle      *webui.WebUIBundle
}

func (app *CovenApplication) Init() bool {
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
			"default config file path", DefaultConfigPath,
		)
	}

	coreConfig := loadConfigFromFile(configPath)
	if coreConfig == nil {
		return false
	}

	if !verifyCoreConfig(coreConfig) {
		return false
	}

	middlewareConfig := loadMiddlewareConfig(coreConfig.MiddlewaresConfigPath)
	if coreConfig.WebUIConfigPath != "" {
		webuiConfig := loadWebUIConfig(coreConfig.WebUIConfigPath)
		if webuiConfig == nil {
			slog.Error("failed to load web UI configuration")
			return false
		}
		err := validateConfig(webuiConfig)
		if err != nil {
			slog.Error("failed to initialize app", "error message", err.Error())
			return false
		}
		app.webUIConfig = webuiConfig
	} else {
		slog.Warn("web ui config was not provided running bare server with no UI")
	}

	shareFileConfig := loadShareFileConfig(coreConfig.FileServerConfigPath)
	err := validateFileServerConfig(shareFileConfig)
	if err != nil {
		slog.Error("failed to load share file config", "error message", err.Error())
		return false
	}

	app.configuration = coreConfig
	app.middlewareConfig = middlewareConfig
	app.sharFileConfig = shareFileConfig

	configuredLogHandler := configureLogging(coreConfig.LogLevel, coreConfig.LogDirPath)
	completeLogger := slog.New(configuredLogHandler)
	slog.SetDefault(completeLogger)
	return true
}

func (app *CovenApplication) Run() {
	if app.configuration == nil {
		panic("trying to run app with nil configuration")
	}

	fullAddr := fmt.Sprintf("%s:%d", app.configuration.Address,
		app.configuration.Port,
	)

	router := http.NewServeMux()
	mux, err := registerMiddlewares(router, app.middlewareConfig)
	if err != nil {
		slog.Error("failed to configure middlewares")
		return
	}

	if app.webUIConfig != nil {
		err = registerFromEndpoints(router)
		if err != nil {
			slog.Error("failed to register form endpoints")
			return
		}
		uibundle, err := webui.New(app.webUIConfig)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		err = registerWebUIStaticFilesEndpoints(router, app.webUIConfig)
		if err != nil {
			slog.Error("failed to register web ui static files endpoints", "error message", err.Error())
			return
		}
		err = registerUIEndpoints(router, uibundle)
		if err != nil {
			slog.Error("failed to regeister ui endpoints", "error message", err.Error())
			return
		}
		app.webUIBundle = uibundle
	} else {
		slog.Warn("web ui bundle config was not provided form and UI endpoints are disabled")
	}

	if app.sharFileConfig != nil {
		err = registerFileShareEndpoints(router, app.sharFileConfig)
		if err != nil {
			slog.Error("failed to register files endpoints", "error message", err.Error())
			return
		}
		slog.Info("succesfly init file server")
	} else {
		slog.Info("shared files are disabled")
	}

	serv := http.Server{
		Addr:         fullAddr,
		ReadTimeout:  app.configuration.ReadTimeout * time.Second,
		WriteTimeout: app.configuration.WriteTimeout * time.Second,
		Handler:      mux,
	}

	isHttps := app.configuration.CertFilePath != ""

	if isHttps {
		slog.Info(fmt.Sprintf("running service at https://%s", fullAddr))
		err := serv.ListenAndServeTLS(app.configuration.CertFilePath,
			app.configuration.CertKeyFilePath,
		)
		if err != nil {
			slog.Error("failed to run server", "reason", err.Error())
		}
	} else {
		slog.Info(fmt.Sprintf("running service at http://%s", fullAddr))
		err := serv.ListenAndServe()
		if err != nil {
			slog.Error("failed to run server", "reason", err.Error())
		}
	}
}
