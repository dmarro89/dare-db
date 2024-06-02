package server

import (
	"os"
	"testing"

	"github.com/dmarro89/dare-db/database"
	"gotest.tools/assert"
)

// TestMain runs setup before tests and teardown after tests
func TestMain(m *testing.M) {

	// Init configuration first
	SetupTestConfiguration()

	// Run the tests
	code := m.Run()

	// Teardown code here
	os.Unsetenv("DARE_TLS_ENABLED")
	TeardownTestConfiguration()

	// Exit with the proper code
	os.Exit(code)
}

func TestNewServerFactory(t *testing.T) {
	factory := NewFactory()
	assert.Assert(t, factory != nil, "NewServerFactory() should not return nil")
}

func TestNewServerWithTlsEnabled(t *testing.T) {
	// Set the DARE_TLS_ENABLED environment variable to "true"
	t.Setenv("DARE_TLS_ENABLED", "true")

	factory := NewFactory()
	server := factory.GetWebServer(NewDareServer(database.NewDatabase()))

	// Assert that the server is of type HttpsServer
	_, isHttpsServer := server.(*HttpsServer)
	assert.Assert(t, !isHttpsServer, "NewServer() should return an HttpsServer when TLS_ENABLED is false")
}

func TestNewServerWithTlsDisabled(t *testing.T) {
	// Set the DARE_TLS_ENABLED environment variable to "false"
	t.Setenv("DARE_TLS_ENABLED", "false")

	factory := NewFactory()
	server := factory.GetWebServer(NewDareServer(database.NewDatabase()))

	// Assert that the server is of type HttpServer
	_, isHttpServer := server.(*HttpServer)
	assert.Assert(t, isHttpServer, "NewServer() should return an HttpServer when TLS_ENABLED is false")
}

func TestGetTlsEnabled(t *testing.T) {
	factory := NewFactory()

	// Test when DARE_TLS_ENABLED is "true"
	t.Setenv("DARE_TLS_ENABLED", "true")
	reReadConfigsFromEnvs()
	assert.Assert(t, factory.getTLSEnabled(), "getTLSEnabled() should return true when TLS_ENABLED is 'true'")

	// Test when DARE_TLS_ENABLED is "false"
	t.Setenv("DARE_TLS_ENABLED", "false")
	reReadConfigsFromEnvs()
	assert.Assert(t, !factory.getTLSEnabled(), "getTLSEnabled() should return false when TLS_ENABLED is 'false'")

	// Test when DARE_TLS_ENABLED is not set
	os.Unsetenv("DARE_TLS_ENABLED")
	assert.Assert(t, !factory.getTLSEnabled(), "getTLSEnabled() should return false when TLS_ENABLED is not set")
}
