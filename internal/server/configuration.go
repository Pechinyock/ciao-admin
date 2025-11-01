package server

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
	MiddlewaresConfigPath string        `json:"middlewaresConfigPath"`
}

type MeddlewaresConfig struct {
	RequestLogger *RequestLoggerMiddlewareConfig `json:"requestLogger"`
}

type RequestLoggerMiddlewareConfig struct {
	Enabled  bool   `json:"enabled"`
	LogLevel string `json:"logLevel"`
}
