package server

import (
	"fmt"
	"os"
	"testing"

	"github.com/dmarro89/dare-db/auth"
	"github.com/dmarro89/dare-db/database"
	"github.com/dmarro89/dare-db/logger"
	"gotest.tools/assert"
)

// TestMain runs setup before tests and teardown after tests
func TestMain(m *testing.M) {

	// Init configuration first
	testConf := SetupTestConfiguration()
	fmt.Println("Test log file should be: ", testConf.Get("log.log_file"))

	// Run the tests
	code := m.Run()

	// Teardown code here
	os.Unsetenv("DARE_TLS_ENABLED")
	TeardownTestConfiguration()

	// Exit with the proper code
	os.Exit(code)
}

func TestNewServerFactory(t *testing.T) {
	factory := NewFactory(NewConfiguration(""), logger.NewDareLogger())
	assert.Assert(t, factory != nil, "NewServerFactory() should not return nil")
}

func TestNewServerWithTlsEnabled(t *testing.T) {
	// Set the DARE_TLS_ENABLED environment variable to "true"
	t.Setenv("DARE_TLS_ENABLED", "true")

	factory := NewFactory(NewConfiguration(""), logger.NewDareLogger())
	server := factory.GetWebServer(NewDareServer(database.NewDatabase(), auth.NewUserStore()))

	// Assert that the server is of type HttpsServer
	_, isHttpsServer := server.(*HttpsServer)
	assert.Assert(t, isHttpsServer, "NewServer() should return an HttpsServer when TLS_ENABLED is true")
}

func TestNewServerWithTlsDisabled(t *testing.T) {
	// Set the DARE_TLS_ENABLED environment variable to "false"
	t.Setenv("DARE_TLS_ENABLED", "false")

	factory := NewFactory(NewConfiguration(""), logger.NewDareLogger())
	server := factory.GetWebServer(NewDareServer(database.NewDatabase(), auth.NewUserStore()))

	// Assert that the server is of type HttpServer
	_, isHttpServer := server.(*HttpServer)
	assert.Assert(t, isHttpServer, "NewServer() should return an HttpServer when TLS_ENABLED is false")
}

/*
//FIXME: pass teh right config to the factory
func TestGetTlsEnabled(t *testing.T) {
	factory := NewFactory()

	// Test when DARE_TLS_ENABLED is "true"
	t.Setenv("DARE_TLS_ENABLED", "true")
	//reReadConfigsFromEnvs(factory)
	assert.Assert(t, factory.getTLSEnabled(), "getTLSEnabled() should return true when TLS_ENABLED is 'true'")

	// Test when DARE_TLS_ENABLED is "false"
	t.Setenv("DARE_TLS_ENABLED", "false")
	//reReadConfigsFromEnvs()
	assert.Assert(t, !factory.getTLSEnabled(), "getTLSEnabled() should return false when TLS_ENABLED is 'false'")

	// Test when DARE_TLS_ENABLED is not set
	os.Unsetenv("DARE_TLS_ENABLED")
	assert.Assert(t, !factory.getTLSEnabled(), "getTLSEnabled() should return false when TLS_ENABLED is not set")
}
*/
