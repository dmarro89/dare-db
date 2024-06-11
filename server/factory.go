package server

import (
	"github.com/spf13/viper"
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
	return viper.GetBool("security.tls_enabled")
}
