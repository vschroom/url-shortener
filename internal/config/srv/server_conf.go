package srv

import (
	"flag"

	"os"
)

type ServerConfig struct {
	Addr            string
	BaseShortAddr   string
	LoggerLevel     string
	FileStoragePath string
}

func InitServerConfig() ServerConfig {
	var cfg ServerConfig
	flag.StringVar(&cfg.Addr, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.BaseShortAddr, "b", "http://localhost:8080/", "base address and port for short url")
	flag.StringVar(&cfg.LoggerLevel, "l", "Info", "logger level")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "file storage path")
	flag.Parse()

	if srvAddress := os.Getenv("SERVER_ADDRESS"); srvAddress != "" {
		cfg.Addr = srvAddress
	}

	if baseUrl := os.Getenv("BASE_URL"); baseUrl != "" {
		cfg.BaseShortAddr = baseUrl
	}

	if loggerLevel := os.Getenv("LOGGER_LEVEL"); loggerLevel != "" {
		cfg.LoggerLevel = loggerLevel
	}

	if fileStoragePath := os.Getenv("FILE_STORAGE_PATH"); fileStoragePath != "" {
		cfg.FileStoragePath = fileStoragePath
	}

	return cfg
}
