package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmarro89/dare-db/logger"
)

type Server interface {
	Start()
	Stop()
}

type HttpServer struct {
	dareServer    IDare
	httpServer    *http.Server
	configuration Config
	sigChan       chan os.Signal
	logger        logger.Logger
}

func NewHttpServer(dareServer IDare, configuration Config, logger logger.Logger) *HttpServer {
	return &HttpServer{
		dareServer:    dareServer,
		configuration: configuration,
		sigChan:       make(chan os.Signal, 1),
		logger:        logger,
	}
}

func (server *HttpServer) Start() {
	if server.configuration.IsSet("log.log_file") {
		server.logger.Start(server.configuration.GetString("log.log_file"))
	}

	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", server.configuration.GetString("server.host"), server.configuration.GetString("server.port")),
		Handler: server.dareServer.CreateMux(nil, nil),
	}

	go func() {
		server.logger.Info("Serving new connections on: ", server.configuration.GetString("server.host"), ":", server.configuration.GetString("server.port"))
		if err := server.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			server.logger.Fatal("HTTP server error: %v", err)
		}
		server.logger.Info("Stopped serving new connections.")
	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.httpServer.Shutdown(shutdownCtx); err != nil {
		server.logger.Fatal("HTTP shutdown error:", err)
	}

	server.logger.Info("Graceful shutdown complete.")
	server.httpServer = nil

	server.logger.Close()
}

type HttpsServer struct {
	dareServer    IDare
	httpsServer   *http.Server
	configuration Config
	sigChan       chan os.Signal
	logger        logger.Logger
}

func NewHttpsServer(dareServer IDare, configuration Config, logger logger.Logger) *HttpsServer {
	return &HttpsServer{
		sigChan:       make(chan os.Signal, 1),
		configuration: configuration,
		dareServer:    dareServer,
		logger:        logger,
	}
}

func (server *HttpsServer) Start() {
	server.httpsServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", server.configuration.GetString("server.host"), server.configuration.GetString("server.port")),
		Handler: server.dareServer.CreateMux(nil, nil),
	}

	go func() {
		server.logger.Info("Serving new connections on: ", server.configuration.GetString("server.host"), ":", server.configuration.GetString("server.port"))
		server.logger.Info("Using certificate files. (1) ", server.configuration.GetString("security.cert_private"), " ; (2) ", server.configuration.GetString("security.cert_public"))

		if err := server.httpsServer.ListenAndServeTLS(server.configuration.GetString("security.cert_public"), server.configuration.GetString("security.cert_private")); !errors.Is(err, http.ErrServerClosed) {
			server.logger.Fatal("HTTPS server error: ", err)
		}
		server.logger.Info("Stopped serving new connections.")
		server.logger.Close()
	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpsServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.httpsServer.Shutdown(shutdownCtx); err != nil {
		server.logger.Fatal("HTTP shutdown error:", err)
	}

	server.logger.Info("Graceful shutdown complete.")
	server.httpsServer = nil
	server.logger.Close()
}
