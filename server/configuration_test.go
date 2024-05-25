package server

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestNewConfiguration(t *testing.T) {
	// Set environment variables
	os.Setenv("DARE_PORT", "8080")
	os.Setenv("DARE_HOST", "localhost")
	os.Setenv("TLS_ENABLED", "true")
	os.Setenv("TLS_CERT_FILE", "/path/to/cert")
	os.Setenv("TLS_KEY_FILE", "/path/to/key")

	// Create a new configuration
	config := NewConfiguration()

	// Check if the values are correctly set
	assert.Equal(t, "8080", config.Port, "Port should be '8080'")
	assert.Equal(t, "localhost", config.Host, "Host should be 'localhost'")
	assert.Equal(t, config.TLSEnabled, true, "TLSEnabled should be true")
	assert.Equal(t, "/path/to/cert", config.TLSCertFile, "TLSCertFile should be '/path/to/cert'")
	assert.Equal(t, "/path/to/key", config.TLSKeyFile, "TLSKeyFile should be '/path/to/key'")

	// Cleanup environment variables
	os.Unsetenv("DARE_PORT")
	os.Unsetenv("DARE_HOST")
	os.Unsetenv("TLS_ENABLED")
	os.Unsetenv("TLS_CERT_FILE")
	os.Unsetenv("TLS_KEY_FILE")
}

func TestNewConfigurationDefaults(t *testing.T) {
	// Ensure environment variables are not set
	os.Unsetenv("DARE_PORT")
	os.Unsetenv("DARE_HOST")
	os.Unsetenv("TLS_ENABLED")
	os.Unsetenv("TLS_CERT_FILE")
	os.Unsetenv("TLS_KEY_FILE")

	// Create a new configuration
	config := NewConfiguration()

	// Check if the values are correctly set to defaults
	assert.Equal(t, "2605", config.Port, "Port should be '2605'")
	assert.Equal(t, "", config.Host, "Host should be ''")
	assert.Equal(t, config.TLSEnabled, false, "TLSEnabled should be false")
	assert.Equal(t, "", config.TLSCertFile, "TLSCertFile should be ''")
	assert.Equal(t, "", config.TLSKeyFile, "TLSKeyFile should be ''")
}

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("TEST_ENV", "value")

	result := getEnvOrDefault("TEST_ENV", "default")
	assert.Equal(t, "value", result, "Result should be 'value'")

	result = getEnvOrDefault("MISSING_ENV", "default")
	assert.Equal(t, "default", result, "Result should be 'default'")

	os.Unsetenv("TEST_ENV")
}

func TestGetEnvBooleanOrDefault(t *testing.T) {
	os.Setenv("BOOL_ENV", "true")

	result := getEnvBooleanOrDefault("BOOL_ENV", false)
	assert.Equal(t, result, true, "Result should be true")

	result = getEnvBooleanOrDefault("MISSING_BOOL_ENV", false)
	assert.Equal(t, result, false, "Result should be false")

	os.Unsetenv("BOOL_ENV")
}
