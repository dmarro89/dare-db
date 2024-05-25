package server

import (
	"os"
	"strings"
)

type ServerFactory struct {
}

func NewServerFactory() *ServerFactory {
	return &ServerFactory{}
}

func (f *ServerFactory) NewServer() Server {
	if f.getTLSEnabled() {
		return NewHttpsServer()
	}

	return NewHttpServer()
}

func (f *ServerFactory) getTLSEnabled() bool {
	tlsEnabled := os.Getenv("TLS_ENABLED")
	return strings.EqualFold(tlsEnabled, "true")
}
