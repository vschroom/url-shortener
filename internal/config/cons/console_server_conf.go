package cons

import (
	"flag"
)

type ServerConsoleArg struct {
	Addr          string
	BaseShortAddr string
}

func ParseServerFlags() ServerConsoleArg {
	srvConsArgs := ServerConsoleArg{}
	flag.StringVar(&srvConsArgs.Addr, "a", ":8080", "address and port to run server")
	flag.StringVar(&srvConsArgs.BaseShortAddr, "b", "http://localhost:8080/", "base address and port for short url")

	flag.Parse()

	return srvConsArgs
}
