package cons

import (
	"flag"
)

var ServerConsoleArg struct {
	Addr          string
	BaseShortAddr string
}

func ParseServerFlags() {
	flag.StringVar(&ServerConsoleArg.Addr, "a", ":8080", "address and port to run server")
	flag.StringVar(&ServerConsoleArg.BaseShortAddr, "b", "localhost:8080", "base address and port for short url")

	flag.Parse()
}
