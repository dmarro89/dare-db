package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-while/nodare-db/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Start()
	Stop()
}

type HttpServer struct {
	ndbServer     WebMux
	httpServer    *http.Server
	configuration Config
	sigChan       chan os.Signal
	logger        *ilog.LOG
}

type HttpsServer struct {
	ndbServer     WebMux
	httpsServer   *http.Server
	configuration Config
	sigChan       chan os.Signal
	logger        *ilog.LOG
}

func NewHttpServer(ndbServer WebMux, logger *ilog.LOG) (srv *HttpServer, sub_dicks uint32) {
	cfg, sub_dicks := NewConfiguration("")
	srv = &HttpServer{
		ndbServer:     ndbServer,
		configuration: cfg,
		sigChan:       make(chan os.Signal, 1),
		logger:        logger,
	}
	return
}

func NewHttpsServer(ndbServer WebMux, logger *ilog.LOG) (srv *HttpsServer, sub_dicks uint32) {
	cfg, sub_dicks := NewConfiguration("")
	srv = &HttpsServer{
		sigChan:       make(chan os.Signal, 1),
		configuration: cfg,
		ndbServer:     ndbServer,
		logger:        logger,
	}
	return
}

func (server *HttpServer) Start() {

	if server.configuration.IsSet("log.log_file") {
		server.logger.OpenLogFile(server.configuration.GetString("log.log_file"))
	}

	server.httpServer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         fmt.Sprintf("%s:%s", server.configuration.GetString("server.host"), server.configuration.GetString("server.port")),
		Handler:      server.ndbServer.CreateMux(),
	}

	go func() {
		server.logger.Info("HttpServer @ '%s:%s'", server.configuration.GetString("server.host"), server.configuration.GetString("server.port"))
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
		Addr:         fmt.Sprintf("%s:%s", server.configuration.GetString("server.host"), server.configuration.GetString("server.port")),
		Handler:      server.ndbServer.CreateMux(),
	}

	go func() {
		server.logger.Info("HttpsServer @ '%s:%s'", server.configuration.GetString("server.host"), server.configuration.GetString("server.port"))
		server.logger.Debug("HttpsServer: PUB_CERT='%s' PRIV_KEY='%s'", server.configuration.GetString("security.tls_cert_public"), server.configuration.GetString("security.tls_cert_private"))

		if err := server.httpsServer.ListenAndServeTLS(server.configuration.GetString("security.tls_cert_public"), server.configuration.GetString("security.tls_cert_private")); !errors.Is(err, http.ErrServerClosed) {
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
