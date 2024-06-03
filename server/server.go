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
}

func NewHttpServer(dareServer IDare) *HttpServer {
	return &HttpServer{
		configuration: NewConfiguration(""),
		sigChan:       make(chan os.Signal, 1),
		dareServer:    dareServer,
	}
}

func (server *HttpServer) Start() {

	if server.configuration.IsSet("log.log_file") {
		logger.OpenLogFile(server.configuration.GetString("log.log_file"))
	}

	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", server.configuration.GetString("server.host"), server.configuration.GetString("server.port")),
		Handler: server.dareServer.CreateMux(),
	}

	go func() {
		logger.Info("Serving new connections on: ", server.configuration.GetString("server.host"), ":", server.configuration.GetString("server.port"))
		if err := server.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("HTTP server error: %v", err)
		}
		logger.Info("Stopped serving new connections.")
		logger.CloseLogFile()
	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("HTTP shutdown error:", err)
	}

	logger.Info("Graceful shutdown complete.")
	server.httpServer = nil

	logger.CloseLogFile()
}

type HttpsServer struct {
	dareServer    IDare
	httpsServer   *http.Server
	configuration Config
	sigChan       chan os.Signal
}

func NewHttpsServer(dareServer IDare) *HttpsServer {
	return &HttpsServer{
		configuration: NewConfiguration(""),
		sigChan:       make(chan os.Signal, 1),
		dareServer:    dareServer,
	}
}

func (server *HttpsServer) Start() {
	server.httpsServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", server.configuration.GetString("server.host"), server.configuration.GetString("server.port")),
		Handler: server.dareServer.CreateMux(),
	}

	go func() {
		logger.Info("Serving new connections on: ", server.configuration.GetString("server.host"), ":", server.configuration.GetString("server.port"))
		if err := server.httpsServer.ListenAndServeTLS(server.configuration.GetString("security.tls_cert_private"), server.configuration.GetString("security.tls_cert_public")); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("HTTPS server error: ", err)
		}
		logger.Info("Stopped serving new connections.")
		logger.CloseLogFile()
	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpsServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.httpsServer.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("HTTP shutdown error:", err)
	}

	logger.Info("Graceful shutdown complete.")
	server.httpsServer = nil
	logger.CloseLogFile()
}
