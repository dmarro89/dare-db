package server

import (
	"os"
	"strconv"
	"github.com/dmarro89/dare-db/logger"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) GetWebServer(dareServer IDare, logger *darelog.LOG) Server {

	if f.getTLSEnabled() {
		return NewHttpsServer(dareServer, logger)
	}

	return NewHttpServer(dareServer, logger)
}

func (f *Factory) getTLSEnabled() bool {
	//FIXME: pass teh right config to the factory
	isTLSEnabled, err := strconv.ParseBool(os.Getenv("DARE_TLS_ENABLED"))
	if err != nil {
		isTLSEnabled = false
	}
	return isTLSEnabled
}
