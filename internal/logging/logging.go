package logging

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func GetOrCreateLogFile(logDirPath string) (*os.File, error) {
	if logDirPath == "" {
		slog.Error("failed to create log file, provided log direcotry path is empty")
		return nil, errors.New("log directory path is empty")
	}

	var fullPath string
	if !filepath.IsAbs(logDirPath) {
		tryFullPath, err := filepath.Abs(logDirPath)
		if err != nil {
			slog.Error("failed to convert file path to absolute",
				"error message", err.Error(),
				"provided path", logDirPath,
			)
			return nil, err
		}
		fullPath = tryFullPath
	} else {
		fullPath = logDirPath
	}

	info, err := os.Stat(fullPath)
	if err != nil && !os.IsNotExist(err) {

		slog.Error("failed to get information about direcotry",
			"error message", err.Error(),
			"path to directory", fullPath,
		)
		return nil, err
	}

	if info == nil {
		slog.Warn("log directory is not exist, trying to create it...")

		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			slog.Error("failed to create log directory",
				"error message", err.Error(),
				"provided path", logDirPath,
			)
			return nil, err
		} else {
			slog.Info("log directory has been created succesfly",
				"directory path", logDirPath,
			)
		}
	} else if !info.IsDir() {
		msg := "failed to set log directory: provided path is not a path to directory"
		slog.Error(msg)
		return nil, errors.New(msg)
	}
	logFileName := time.Now().Format("2006-01-02") + ".log"
	logFileFullPath := filepath.Join(logDirPath, logFileName)
	if _, err := os.Stat(logFileFullPath); err == nil {
		existingFile, err := os.Open(logFileFullPath)
		if err != nil {
			slog.Error("failed to open existing log file",
				"error message", err.Error(),
				"log file path", logFileFullPath,
			)
			return nil, err
		}
		return existingFile, nil
	}
	logFile, err := os.Create(logFileFullPath)
	if err != nil {
		slog.Error("failed to create log file",
			"error message", err.Error(),
			"log dir path", logDirPath,
		)
		return nil, err
	}
	return logFile, nil
}
