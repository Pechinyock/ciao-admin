package app

import (
	"ciao-admin/internal/logging"
	"ciao-admin/internal/loghandlers"
	"ciao-admin/internal/server"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const DefaultConfigPath = "./configs/debug-server-config.json"

func getFullPath(path string) string {
	var fullFilePath string
	if filepath.IsAbs(path) {
		fullFilePath = path
	} else {
		full, err := filepath.Abs(path)
		if err != nil {
			slog.Error("an error ocured while trying to convert relative path into full path", "error message", err.Error())
			panic("")
		}
		fullFilePath = full
	}
	return fullFilePath
}

func loadConfigFromFile(path string) *server.CoreConfig {
	if path == "" {
		slog.Error("failed to load configuration path to configuration file is empty")
		return nil
	}

	var fullFilePath = getFullPath(path)

	if _, err := os.Stat(fullFilePath); err != nil {
		slog.Error("configuration file doesn't exist", "provided file path", fullFilePath)
		return nil
	}
	serverConfigData, err := os.ReadFile(fullFilePath)
	if err != nil {
		slog.Error("an error occured while trying to load configuration data", "error message", err.Error())
		return nil
	}
	var serverConfig server.CoreConfig
	err = json.Unmarshal(serverConfigData, &serverConfig)
	if err != nil {
		slog.Error("an error occured while trying to parse configuration file", "error message", err.Error())
		return nil
	}
	return &serverConfig
}

func verifyServerConfig(config *server.CoreConfig) bool {
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
		certPath := getFullPath(cert)
		if _, err := os.Stat(certPath); err != nil {
			slog.Error("failed to load certificate file", "file path", certPath)
			result = false
		}

		certKeyPath := getFullPath(certKey)
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
