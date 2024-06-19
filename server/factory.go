package server

import (
	"github.com/dmarro89/dare-db/logger"
)

type Factory struct {
	configuration Config
	logger        logger.Logger
}

func NewFactory(configuration Config, logger logger.Logger) *Factory {
	return &Factory{configuration: configuration, logger: logger}
}

func (f *Factory) GetWebServer(dareServer IDare) Server {
	if f.configuration.GetBool("security.tls_enabled") {
		return NewHttpsServer(dareServer, f.configuration, f.logger)
	}

	return NewHttpServer(dareServer, f.configuration, f.logger)
}
