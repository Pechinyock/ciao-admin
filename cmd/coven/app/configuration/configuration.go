package configuration

import (
	"time"
)

type CoreConfig struct {
	Address               string        `json:"address"`
	Port                  uint16        `json:"port"`
	LogLevel              string        `json:"logLevel"`
	LogDirPath            string        `json:"logDirPath"`
	ReadTimeout           time.Duration `json:"readTimeout"`
	WriteTimeout          time.Duration `json:"writeTimeout"`
	CertFilePath          string        `json:"certFilePath"`
	CertKeyFilePath       string        `json:"certKeyFilePath"`
	EnableApiEndpoints    bool          `json:"enableApiEndpoints"`
	MiddlewaresConfigPath string        `json:"middlewaresConfigPath"`
	WebUIConfigPath       string        `json:"webUIConfigPath"`
	FileServerConfigPath  string        `json:"fileServerConfigPath"`
}

type MeddlewaresConfig struct {
	RequestLogger *RequestLoggerMiddlewareConfig `json:"requestLogger"`
}

type RequestLoggerMiddlewareConfig struct {
	Enabled  bool   `json:"enabled"`
	LogLevel string `json:"logLevel"`
}

type ShareDirConfig struct {
	RouteName   string `json:"routeName"`
	DirPath     string `json:"dirPath"`
	TokenSource string `jsong:"tokenSource"` // header, cookie, none
}

type FileServerConfig struct {
	ShareDirConfigs []ShareDirConfig `json:"shareDirConfigs"`
}
