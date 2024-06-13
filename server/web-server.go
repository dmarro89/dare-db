package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-while/nodare-db/logger"
	//"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Start()
	Stop()
	//Config()
}

type HttpServer struct {
	ndbServer     WebMux
	httpServer    *http.Server
	VCFG           VConfig
	sigChan       chan os.Signal
	logger        *ilog.LOG
}

type HttpsServer struct {
	ndbServer     WebMux
	httpsServer   *http.Server
	VCFG          VConfig
	sigChan       chan os.Signal
	logger        *ilog.LOG
}

func NewHttpServer(ndbServer WebMux, logger *ilog.LOG) (srv *HttpServer, cfg VConfig, sub_dicks uint32) {
	cfg, sub_dicks = NewConfiguration("")
	//log.Printf("NewHttpServer cfg='%#v' ViperConfig='%#v'", cfg, cfg.ViperConfig.)
	srv = &HttpServer{
		sigChan:       make(chan os.Signal, 1),
		ndbServer:     ndbServer,
		logger:        logger,
		VCFG:          cfg,
	}
	return
}

func NewHttpsServer(ndbServer WebMux, logger *ilog.LOG) (srv *HttpsServer, cfg VConfig, sub_dicks uint32) {
	cfg, sub_dicks = NewConfiguration("")
	srv = &HttpsServer{
		sigChan:       make(chan os.Signal, 1),
		ndbServer:     ndbServer,
		logger:        logger,
		VCFG:          cfg,
	}
	return
}

func (server *HttpServer) Start() {

	if server.VCFG.IsSet("log.log_file") {
		server.logger.OpenLogFile(server.VCFG.GetString("log.log_file"))
	}

	server.httpServer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         fmt.Sprintf("%s:%s", server.VCFG.GetString("server.host"), server.VCFG.GetString("server.port")),
		Handler:      server.ndbServer.CreateMux(),
	}

	go func() {
		server.logger.Info("HttpServer @ '%s:%s'", server.VCFG.GetString("server.host"), server.VCFG.GetString("server.port"))
		if err := server.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			server.logger.Fatal("HTTP server error: %v", err)
		}
		server.logger.Info("HttpServer: closing")
		server.logger.CloseLogFile()
	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.httpServer.Shutdown(shutdownCtx); err != nil {
		server.logger.Fatal("HttpServer: shutdown error %v", err)
	}

	server.logger.Info("HttpServer shutdown complete")
	server.httpServer = nil

	server.logger.CloseLogFile()
}

func (server *HttpsServer) Start() {
	server.httpsServer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         fmt.Sprintf("%s:%s", server.VCFG.GetString("server.host"), server.VCFG.GetString("server.port")),
		Handler:      server.ndbServer.CreateMux(),
	}

	go func() {
		server.logger.Info("HttpsServer @ '%s:%s'", server.VCFG.GetString("server.host"), server.VCFG.GetString("server.port"))
		server.logger.Debug("HttpsServer: PUB_CERT='%s' PRIV_KEY='%s'", server.VCFG.GetString("security.tls_cert_public"), server.VCFG.GetString("security.tls_cert_private"))

		if err := server.httpsServer.ListenAndServeTLS(server.VCFG.GetString("security.tls_cert_public"), server.VCFG.GetString("security.tls_cert_private")); !errors.Is(err, http.ErrServerClosed) {
			server.logger.Fatal("HttpsServer: error %v", err)
		}
		server.logger.Debug("HttpsServer: closing")
		server.logger.CloseLogFile()
	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpsServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.httpsServer.Shutdown(shutdownCtx); err != nil {
		server.logger.Fatal("HttpsServer: shutdown error %v", err)
	}

	server.logger.Info("HttpsServer: shutdown complete")
	server.httpsServer = nil
	server.logger.CloseLogFile()
}
