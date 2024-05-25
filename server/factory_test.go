package server

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

// TestMain runs setup before tests and teardown after tests
func TestMain(m *testing.M) {
	// Run the tests
	code := m.Run()

	// Teardown code here
	os.Unsetenv("TLS_ENABLED")

	// Exit with the proper code
	os.Exit(code)
}

func TestNewServerFactory(t *testing.T) {
	factory := NewServerFactory()
	assert.Assert(t, factory != nil, "NewServerFactory() should not return nil")
}

func TestNewServerWithTlsEnabled(t *testing.T) {
	// Set the TLS_ENABLED environment variable to "true"
	os.Setenv("TLS_ENABLED", "true")

	factory := NewServerFactory()
	server := factory.NewServer()

	// Assert that the server is of type HttpsServer
	_, isHttpsServer := server.(*HttpsServer)
	assert.Assert(t, isHttpsServer, "NewServer() should return an HttpsServer when TLS_ENABLED is true")

	// Cleanup
	os.Unsetenv("TLS_ENABLED")
}

func TestNewServerWithTlsDisabled(t *testing.T) {
	// Set the TLS_ENABLED environment variable to "false"
	os.Setenv("TLS_ENABLED", "false")

	factory := NewServerFactory()
	server := factory.NewServer()

	// Assert that the server is of type HttpServer
	_, isHttpServer := server.(*HttpServer)
	assert.Assert(t, isHttpServer, "NewServer() should return an HttpServer when TLS_ENABLED is false")

	// Cleanup
	os.Unsetenv("TLS_ENABLED")
}

func TestGetTlsEnabled(t *testing.T) {
	factory := NewServerFactory()

	// Test when TLS_ENABLED is "true"
	os.Setenv("TLS_ENABLED", "true")
	assert.Assert(t, factory.getTLSEnabled(), "getTLSEnabled() should return true when TLS_ENABLED is 'true'")

	// Test when TLS_ENABLED is "false"
	os.Setenv("TLS_ENABLED", "false")
	assert.Assert(t, !factory.getTLSEnabled(), "getTLSEnabled() should return false when TLS_ENABLED is 'false'")

	// Test when TLS_ENABLED is not set
	os.Unsetenv("TLS_ENABLED")
	assert.Assert(t, !factory.getTLSEnabled(), "getTLSEnabled() should return false when TLS_ENABLED is not set")
}
