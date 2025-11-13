package app

import (
	"ciao-admin/cmd/coven/app/configuration"
	"ciao-admin/internal/UI/webui"
	"ciao-admin/internal/logging"
	"ciao-admin/internal/loghandlers"
	"ciao-admin/internal/utils"
	"errors"
	"path"

	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

const DefaultConfigPath = "./configs/debug-server-config.json"

func loadConfigFromFile(path string) *configuration.CoreConfig {
	if path == "" {
		slog.Error("failed to load configuration path to configuration file is empty")
		return nil
	}

	var fullFilePath = utils.GetFullPath(path)

	if _, err := os.Stat(fullFilePath); err != nil {
		slog.Error("configuration file doesn't exist", "provided file path", fullFilePath)
		return nil
	}
	serverConfigData, err := os.ReadFile(fullFilePath)
	if err != nil {
		slog.Error("an error occured while trying to load configuration data", "error message", err.Error())
		return nil
	}
	var serverConfig configuration.CoreConfig
	err = json.Unmarshal(serverConfigData, &serverConfig)
	if err != nil {
		slog.Error("an error occured while trying to parse configuration file", "error message", err.Error())
		return nil
	}
	return &serverConfig
}

func verifyCoreConfig(config *configuration.CoreConfig) bool {
	var result bool = true
	if config.Address == "" {
		slog.Error("'Domain' must not be empty")
		result = false
	}

	addr := config.Address
	if ip := net.ParseIP(addr); ip == nil {
		slog.Warn("using not raw ip addres")
		pattern := `^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`

		isValidDomainName, err := regexp.MatchString(pattern, addr)
		if err != nil || !isValidDomainName && addr != "localhost" {
			slog.Error(fmt.Sprintf("domain name %s is not valid", addr),
				"must contain only", "letters, numbers, hyphens, and dots",
			)
			result = false
		}
	}

	var readTimeout, writeTimeout = config.ReadTimeout, config.WriteTimeout
	if readTimeout == 0 {
		slog.Warn("read timeout is disabled", "provided value is", readTimeout)
	}
	if writeTimeout == 0 {
		slog.Warn("write timeout is disabled", "provided value is", writeTimeout)
	}

	maxReasonable := 24 * time.Hour
	if readTimeout > maxReasonable {
		slog.Error(fmt.Sprintf("read timeout %v is too large - maximum recommended is 24h", readTimeout))
		result = false
	}
	if writeTimeout > maxReasonable {
		slog.Error(fmt.Sprintf("write timeout %v is too large - maximum recommended is 24h", writeTimeout))
		result = false
	}

	cert, certKey := config.CertFilePath, config.CertKeyFilePath
	if cert == "" && certKey == "" {
		slog.Warn("TLS certificate not provided - running in HTTP mode (not secure)")
	}

	if (cert == "" && certKey != "") || (cert != "" && certKey == "") {
		slog.Warn("incomplete TLS configuration: both certificate and key files must be provided together. Falling back to HTTP mode")
	}

	if cert != "" && certKey != "" {
		certPath := utils.GetFullPath(cert)
		if _, err := os.Stat(certPath); err != nil {
			slog.Error("failed to load certificate file", "file path", certPath)
			result = false
		}

		certKeyPath := utils.GetFullPath(certKey)
		if _, err := os.Stat(certKeyPath); err != nil {
			slog.Error("failed to load certificate key file", "key file path", certKeyPath)
			result = false
		}
	}

	return result
}

func configureLogging(logLevelStr string, logDirPath string) slog.Handler {
	logLevel := parseLogLevel(logLevelStr)
	prettyStdhandler := loghandlers.PrettyStdLogHandler{
		Level:  logLevel,
		Writer: os.Stdout,
	}

	if logDirPath != "" {
		logFile, err := logging.GetOrCreateLogFile(logDirPath)
		if err != nil {
			slog.Error("an error occured while trying to create log file",
				"error message", err.Error(),
			)
		}

		slog.Info("logger has been initialized",
			"log level", logLevel.String(),
			"log channels", fmt.Sprintf("stdout, %s", logDirPath),
		)

		fileHandler := loghandlers.FileLogHandler{
			Level:  logLevel,
			Writer: logFile,
		}

		multiHandler := loghandlers.NewMultiHandler(&prettyStdhandler, &fileHandler)
		return multiHandler
	}
	slog.Info("logger has been initialized",
		"log level", logLevel.String(),
		"log channels", "stdout",
	)
	return &prettyStdhandler
}

func parseLogLevel(logLevelStr string) slog.Level {
	lower := strings.ToLower(logLevelStr)
	switch lower {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "err", "error":
		return slog.LevelError
	default:
		{
			slog.Error("failed to parse log level",
				"source str", logLevelStr,
				"available options", "(debug, "+"info, "+"warn, "+"error)",
			)
			slog.Warn("log level has been set to the default value 'INFO'")
			return slog.LevelInfo
		}
	}
}

func loadMiddlewareConfig(path string) *configuration.MeddlewaresConfig {
	if path == "" {
		return nil
	}

	var fullFilePath = utils.GetFullPath(path)
	if !utils.IsFileExists(fullFilePath) {
		slog.Error("configuration file doesn't exist", "provided file path", fullFilePath)
		return nil
	}

	serverConfigData, err := os.ReadFile(fullFilePath)
	if err != nil {
		slog.Error("an error occured while trying to load configuration data", "error message", err.Error())
		return nil
	}
	var middlewareConfig configuration.MeddlewaresConfig
	err = json.Unmarshal(serverConfigData, &middlewareConfig)
	if err != nil {
		slog.Error("an error occured while trying to parse configuration file", "error message", err.Error())
		return nil
	}
	return &middlewareConfig

}

func validateConfig(cfg *webui.WebUIBundleConfig) error {
	if cfg == nil {
		return errors.New("configuration is nil")
	}
	if cfg.RootPath == "" {
		return errors.New("web ui bundle root path cannot be empty")
	}
	if cfg.StaticFilesDirName == "" {
		return errors.New("provided static files directory name is empty")
	}

	staticFilesDirFullPath := path.Join(cfg.RootPath, cfg.StaticFilesDirName)
	fullPath := utils.GetFullPath(staticFilesDirFullPath)
	if !utils.IsDirExists(fullPath) {
		return errors.New("provided static files directory path is not exists of it's not a directory")
	}
	if cfg.StaticFilesRootRouteName == "" {
		return errors.New("provided static files route name is empty")
	}
	return nil
}

func loadShareFileConfig(path string) *configuration.FileServerConfig {
	if path == "" {
		return nil
	}

	var fullFilePath = utils.GetFullPath(path)
	if !utils.IsFileExists(fullFilePath) {
		slog.Error("configuration file doesn't exist", "provided file path", fullFilePath)
		return nil
	}

	serverConfigData, err := os.ReadFile(fullFilePath)
	if err != nil {
		slog.Error("an error occured while trying to load configuration data", "error message", err.Error())
		return nil
	}

	var fileServerConfig configuration.FileServerConfig
	err = json.Unmarshal(serverConfigData, &fileServerConfig)
	if err != nil {
		slog.Error("an error occured while trying to parse configuration file", "error message", err.Error())
		return nil
	}

	return &fileServerConfig
}

func validateFileServerConfig(cfg *configuration.FileServerConfig) error {
	if cfg == nil {
		return errors.New("file server configuration is nil")
	}
	if len(cfg.ShareDirConfigs) == 0 {
		return errors.New("file server has no directory(es) to share")
	}
	uniqueNameChecker := make(map[string]bool)
	isValueAlreadyDefined := func(value string) bool {
		_, exists := uniqueNameChecker[value]
		return exists
	}
	for _, c := range cfg.ShareDirConfigs {
		nm := c.RouteName
		if isValueAlreadyDefined(nm) {
			return fmt.Errorf("provided config has more than one route with name: %s", nm)
		}
		uniqueNameChecker[nm] = true

		dr := c.DirPath
		if isValueAlreadyDefined(dr) {
			return fmt.Errorf("provided config has more than one route with directory path: %s", dr)
		}

		uniqueNameChecker[dr] = true
		fullDirPath := utils.GetFullPath(dr)
		if !utils.IsDirExists(fullDirPath) {
			return fmt.Errorf("directory doesn't exists path: %s", fullDirPath)
		}
	}
	return nil
}
