package utils

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const pathGlobSymbols = "*?["

func GetFullPath(path string) string {
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

func IsFileExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil || fileInfo.IsDir() {
		return false
	}
	return true
}

func IsDirExists(path string) bool {
	dirInfo, err := os.Stat(path)
	if err != nil || !dirInfo.IsDir() {
		return false
	}
	return true
}

func IsGlob(path string) bool {
	return strings.ContainsAny(path, pathGlobSymbols)
}

func GetFileName(path string, withExt bool) string {
	fileName := filepath.Base(path)
	if withExt {
		return fileName
	} else {
		return strings.TrimSuffix(fileName, filepath.Ext(fileName))
	}
}

func IsFilePath(path string) bool {
	cleanPath := filepath.Clean(path)

	base := filepath.Base(cleanPath)
	ext := filepath.Ext(base)

	return ext != "" && base != ext
}
