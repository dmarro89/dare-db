package server

import (
	"github.com/go-while/nodare-db/logger"
	"os"
	//"log"
	"strconv"
	"sync"
)

type Factory struct {
	mux sync.Mutex
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewNDBServer(cfg VConfig, ndbServer WebMux, logs ilog.ILOG) (srv Server) {
	f.mux.Lock()
	defer f.mux.Unlock()

	bootsocket := false
	if bootsocket {
		//_ = NewSocketHandler(srv)
		//sockets.Start()
	}
	tls_enabled := cfg.GetBool(VK_SEC_TLS_ENABLED)
	logfile := cfg.GetString(VK_LOG_LOGFILE)
	logs.LogStart(logfile)
	logs.Info("factory: viper cfg loaded tls_enabled=%t logfile='%s'", tls_enabled, logfile)
	switch tls_enabled {
	case false:
		// TCP WEB SERVER
		srv = NewHttpServer(cfg, ndbServer, logs)
		logs.Debug("Factory TCP WEB\n srv='%#v'\n^EOL\n\n cfg='%#v'\n^EOL loglevel=%d\n\n", srv, cfg, logs.GetLOGLEVEL())
	case true:
		// TLS WEB SERVER
		srv = NewHttpsServer(cfg, ndbServer, logs)
		logs.Debug("Factory TLS WEB\n  srv='%#v'\n^EOL\n\n cfg='%#v'\n^EOL loglevel=%d\n\n", cfg, srv, logs.GetLOGLEVEL())
	}
	return
} // end func NewNDBServer

func (f *Factory) getEnvTLSEnabled() bool {
	isTLSEnabled, _ := strconv.ParseBool(os.Getenv("NDB_TLS_ENABLED"))
	if isTLSEnabled {
		return true
	}
	return false
}
