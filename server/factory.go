package server

import (
	"github.com/go-while/nodare-db/logger"
	"log"
	"os"
	"strconv"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) GetWebServer(ndbServer WebMux, logger *ilog.LOG) (srv Server, sub_dicks uint32) {

	if f.getTLSEnabled() {
		srv, sub_dicks = NewHttpsServer(ndbServer, logger)
		return
	}
	srv, sub_dicks = NewHttpServer(ndbServer, logger)
	return
}

func (f *Factory) getTLSEnabled() bool {
	isTLSEnabled, err := strconv.ParseBool(os.Getenv("NDB_TLS_ENABLED"))
	if err != nil {
		isTLSEnabled = false
	}
	log.Printf("isTLSEnabled=%t", isTLSEnabled)
	return isTLSEnabled
}
