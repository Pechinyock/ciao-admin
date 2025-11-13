package app

import (
	"ciao-admin/internal/UI/webui"
	"ciao-admin/internal/utils"
	"encoding/json"
	"log/slog"
	"os"
)

func loadWebUIConfig(path string) *webui.WebUIBundleConfig {
	if path == "" {
		slog.Error("failed to load ui bundle config, provided path is empty")
		return nil
	}
	fullPath := utils.GetFullPath(path)
	if _, err := os.Stat(fullPath); err != nil {
		slog.Error("configuration file doesn't exist", "provided file path", fullPath)
		return nil
	}
	webUIBundleConfig, err := os.ReadFile(fullPath)
	if err != nil {
		slog.Error("an error occured while trying to load configuration data", "error message", err.Error())
		return nil
	}
	var webUIConfig webui.WebUIBundleConfig
	err = json.Unmarshal(webUIBundleConfig, &webUIConfig)
	if err != nil {
		slog.Error("an error occured while trying to parse configuration file", "error message", err.Error())
		return nil
	}
	return &webUIConfig
}
