package server

import (
	"os"
	"strconv"
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
	//FIXME: pass teh right config to the factory
	isTLSEnabled, err := strconv.ParseBool(os.Getenv("DARE_TLS_ENABLED"))
	if err != nil {
		isTLSEnabled = false
	}
	return isTLSEnabled
}
