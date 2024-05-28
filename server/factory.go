package server

import (
	"os"
	"strings"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) GetWebServer(dareServer IDare) Server {
	if f.getTLSEnabled() {
		return NewHttpsServer(dareServer)
	}

	return NewHttpServer(dareServer)
}

func (f *Factory) getTLSEnabled() bool {
	tlsEnabled := os.Getenv(DARE_TLS_ENABLED)
	return strings.EqualFold(tlsEnabled, "true")
}
