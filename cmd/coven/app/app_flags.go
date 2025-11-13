package app

import (
	"flag"
	"fmt"
	"log/slog"
)

const ConfigPathFlagName = "config"

func readFlags() string {
	configPath := flag.String(ConfigPathFlagName,
		DefaultConfigPath,
		"server configuration file path",
	)

	flag.Parse()

	if configPath == nil || *configPath == "" {
		slog.Error(fmt.Sprintf("flag %s was provided, but value is empty", ConfigPathFlagName))
		slog.Warn("tying to load default config", "config path", DefaultConfigPath)
		return DefaultConfigPath
	}
	slog.Info("loading configuration", "configuration file path", *configPath)

	return *configPath
}
