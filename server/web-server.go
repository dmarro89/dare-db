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
}

type HttpServer struct {
	ndbServer  WebMux
	httpServer *http.Server
	cfg        VConfig
	sigChan    chan os.Signal
	logs       ilog.ILOG
}

type HttpsServer struct {
	ndbServer   WebMux
	httpsServer *http.Server
	cfg         VConfig
	sigChan     chan os.Signal
	logs        ilog.ILOG
}

func NewHttpServer(cfg VConfig, ndbServer WebMux, logs ilog.ILOG) (srv *HttpServer) {
	srv = &HttpServer{
		sigChan:   make(chan os.Signal, 1),
		ndbServer: ndbServer,
		logs:      logs,
		cfg:       cfg,
	}
	return
}

func NewHttpsServer(cfg VConfig, ndbServer WebMux, logs ilog.ILOG) (srv *HttpsServer) {
	srv = &HttpsServer{
		sigChan:   make(chan os.Signal, 1),
		ndbServer: ndbServer,
		logs:      logs,
		cfg:       cfg,
	}
	return
}

func (server *HttpServer) Start() {

	server.httpServer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         fmt.Sprintf("%s:%s", server.cfg.GetString(VK_SERVER_HOST), server.cfg.GetString(VK_SERVER_PORT_TCP)),
		Handler:      server.ndbServer.CreateMux(),
	}

	go func() {
		server.logs.Info("HTTP @ '%s:%s'", server.cfg.GetString(VK_SERVER_HOST), server.cfg.GetString(VK_SERVER_PORT_TCP))
		if err := server.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			server.logs.Fatal("HTTP server error: %v", err)
		}
		server.logs.Info("HttpServer: closing")
	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	if err := server.httpServer.Shutdown(shutdownCtx); err != nil {
		server.logs.Fatal("HttpServer: shutdown error %v", err)
	}
	server.logs.Info("HttpServer shutdown complete")
}

func (server *HttpsServer) Start() {

	server.httpsServer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         fmt.Sprintf("%s:%s", server.cfg.GetString(VK_SERVER_HOST), server.cfg.GetString(VK_SERVER_PORT_TCP)),
		Handler:      server.ndbServer.CreateMux(),
	}

	go func() {
		server.logs.Info("HTTPS @ '%s:%s'", server.cfg.GetString(VK_SERVER_HOST), server.cfg.GetString(VK_SERVER_PORT_TCP))
		server.logs.Debug("HttpsServer: PUB_CERT='%s' PRIV_KEY='%s'", server.cfg.GetString(VK_SEC_TLS_PUBCERT), server.cfg.GetString(VK_SEC_TLS_PRIVKEY))

		if err := server.httpsServer.ListenAndServeTLS(server.cfg.GetString(VK_SEC_TLS_PUBCERT), server.cfg.GetString(VK_SEC_TLS_PRIVKEY)); !errors.Is(err, http.ErrServerClosed) {
			server.logs.Fatal("HttpsServer: error %v", err)
		}
		server.logs.Debug("HttpsServer: closing")

	}()

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-server.sigChan
}

func (server *HttpsServer) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	if err := server.httpsServer.Shutdown(shutdownCtx); err != nil {
		server.logs.Fatal("HttpsServer: shutdown error %v", err)
	}
	server.logs.Info("HttpsServer: shutdown complete")
}
