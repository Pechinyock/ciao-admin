package server

import "time"

type ServerConfig struct {
	Address         string        `json:"address"`
	Port            uint16        `json:"port"`
	LogLevel        string        `json:"logLevel"`
	ReadTimeout     time.Duration `json:"readTimeout"`
	WriteTimeout    time.Duration `json:"writeTimeout"`
	CertFilePath    string        `json:"certFilePath"`
	CertKeyFilePath string        `json:"certKeyFilePath"`
}
