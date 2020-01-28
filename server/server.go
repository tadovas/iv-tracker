package server

import (
	"flag"
	"net/http"
)

type Flags struct {
	Address string
}

func RegisterFlags(flags *Flags) {
	flag.StringVar(&flags.Address, "http.address", "localhost:8080", "Bind address of http server")
}

func Setup(flags Flags, handler http.Handler) (http.Server, error) {
	return http.Server{
		Addr:    flags.Address,
		Handler: handler,
	}, nil
}
