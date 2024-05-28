package server

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	dareServer    IDare
	httpServer    *http.Server
	configuration *Configuration
	sigChan       chan os.Signal
}

func NewHttpServer(dareServer IDare) *HttpServer {
	return &HttpServer{
		configuration: NewConfiguration(),
		sigChan:       make(chan os.Signal, 1),
		dareServer:    dareServer,
	}
}

func (server *HttpServer) Start() {
	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", server.configuration.Host, server.configuration.Port),
		Handler: server.dareServer.CreateMux(),
	}

	go func() {
		log.Println("Serving new connections.")
		if err := server.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Println("Graceful shutdown complete.")
	server.httpServer = nil
}

type HttpsServer struct {
	dareServer    IDare
	httpsServer   *http.Server
	configuration *Configuration
	sigChan       chan os.Signal
}

func NewHttpsServer(dareServer IDare) *HttpsServer {
	return &HttpsServer{
		configuration: NewConfiguration(),
		sigChan:       make(chan os.Signal, 1),
		dareServer:    dareServer,
	}
}

func (server *HttpsServer) Start() {
	server.httpsServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", server.configuration.Host, server.configuration.Port),
		Handler: server.dareServer.CreateMux(),
	}

	go func() {
		log.Println("Serving new connections.")
		if err := server.httpsServer.ListenAndServeTLS(server.configuration.TLSCertFile, server.configuration.TLSKeyFile); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpsServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.httpsServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Println("Graceful shutdown complete.")
	server.httpsServer = nil
}
