package srv

import (
	"flag"

	"os"
)

type ServerConfig struct {
	Addr          string `env:"SERVER_ADDRESS"`
	BaseShortAddr string `env:"BASE_URL"`
}

func InitServerConfig() ServerConfig {
	var cfg ServerConfig
	flag.StringVar(&cfg.Addr, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.BaseShortAddr, "b", "http://localhost:8080/", "base address and port for short url")
	flag.Parse()

	if srvAddress := os.Getenv("SERVER_ADDRESS"); srvAddress != "" {
		cfg.Addr = srvAddress
	}

	if baseUrl := os.Getenv("BASE_URL"); baseUrl != "" {
		cfg.BaseShortAddr = baseUrl
	}

	return cfg
}
