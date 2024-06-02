package server

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
)

// constants
const SLEEP_ON_START int = 3

// Mock Database
type MockDatabase struct {
	mock.Mock
}

// Mock DareServer
type MockDareServer struct {
	mock.Mock
}

func (ds *MockDareServer) CreateMux() *http.ServeMux {
	args := ds.Called()
	return args.Get(0).(*http.ServeMux)
}

func (ds *MockDareServer) HandlerGetById(w http.ResponseWriter, r *http.Request) {
}
func (ds *MockDareServer) HandlerSet(w http.ResponseWriter, r *http.Request) {
}
func (ds *MockDareServer) HandlerDelete(w http.ResponseWriter, r *http.Request) {
}

func TestNewHttpServer(t *testing.T) {
	server := NewHttpServer(&MockDareServer{})
	assert.Assert(t, server != nil)
	//assert.Assert(t, server.configuration != nil)
	assert.Assert(t, server.sigChan != nil)
	assert.Assert(t, server.dareServer != nil)
}

func TestHttpServerStartAndStop(t *testing.T) {
	// Setup
	sigChan := make(chan os.Signal, 1)
	server := &HttpServer{
		//configuration: NewConfiguration(),
		sigChan:    sigChan,
		dareServer: &MockDareServer{},
	}

	mux := http.NewServeMux()
	server.dareServer.(*MockDareServer).On("CreateMux").Return(mux)

	// Start server
	go server.Start()
	time.Sleep(time.Duration(SLEEP_ON_START) * time.Second) // Give it time to start

	// Verify the server is running
	assert.Assert(t, server.httpServer != nil)

	// Stop server
	server.Stop()

	time.Sleep(time.Duration(SLEEP_ON_START) * time.Second)
	// Verify the server is stopped
	assert.Assert(t, server.httpServer == nil)
}

func TestNewHttpsServer(t *testing.T) {
	server := NewHttpsServer(&MockDareServer{})
	assert.Assert(t, server != nil)
	//assert.Assert(t, server.configuration != nil)
	assert.Assert(t, server.sigChan != nil)
	assert.Assert(t, server.dareServer != nil)
}

func TestHttpsServerStartAndStop(t *testing.T) {
	t.Skip("Skipping test - Configure certificates to run it")
	// Setup
	sigChan := make(chan os.Signal, 1)
	server := &HttpsServer{
		//configuration: NewConfiguration(),
		sigChan:    sigChan,
		dareServer: &MockDareServer{},
	}

	mux := http.NewServeMux()
	server.dareServer.(*MockDareServer).On("CreateMux").Return(mux)

	// Start server
	go server.Start()
	time.Sleep(time.Duration(SLEEP_ON_START) * time.Second) // Give it time to start

	// Verify the server is running
	assert.Assert(t, server.httpsServer != nil)

	// Stop server
	server.Stop()

	// Verify the server is stopped
	assert.Assert(t, server.httpsServer == nil)
}
