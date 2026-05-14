package srv

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type ServerConfig struct {
	Addr          string `env:"SERVER_ADDRESS"`
	BaseShortAddr string `env:"BASE_URL"`
}

func InitServerConfig() ServerConfig {
	var cfg ServerConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Error while parse Env values")
	}

	if cfg.Addr == "" {
		flag.StringVar(&cfg.Addr, "a", ":8080", "address and port to run server")
	}

	if cfg.BaseShortAddr == "" {
		flag.StringVar(&cfg.BaseShortAddr, "b", "http://localhost:8080/", "base address and port for short url")
	}

	flag.Parse()

	return cfg
}
