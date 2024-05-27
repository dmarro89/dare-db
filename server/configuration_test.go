package server

import (
	"testing"

	"gotest.tools/assert"
)

func TestNewConfiguration(t *testing.T) {
	// Set environment variables
	t.Setenv(DARE_PORT, "8080")
	t.Setenv(DARE_HOST, "localhost")
	t.Setenv(DARE_TLS_ENABLED, "true")
	t.Setenv(DARE_TLS_CERT_FILE, "/path/to/cert")
	t.Setenv(DARE_TLS_KEY_FILE, "/path/to/key")

	// Create a new configuration
	config := NewConfiguration()

	// Check if the values are correctly set
	assert.Equal(t, "8080", config.Port, "Port should be '8080'")
	assert.Equal(t, "localhost", config.Host, "Host should be 'localhost'")
	assert.Equal(t, config.TLSEnabled, true, "TLSEnabled should be true")
	assert.Equal(t, "/path/to/cert", config.TLSCertFile, "TLSCertFile should be '/path/to/cert'")
	assert.Equal(t, "/path/to/key", config.TLSKeyFile, "TLSKeyFile should be '/path/to/key'")
}

func TestNewConfigurationDefaults(t *testing.T) {
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
	t.Setenv("TEST_ENV", "value")

	result := getEnvOrDefault("TEST_ENV", "default")
	assert.Equal(t, "value", result, "Result should be 'value'")

	result = getEnvOrDefault("MISSING_ENV", "default")
	assert.Equal(t, "default", result, "Result should be 'default'")
}

func TestGetEnvBooleanOrDefault(t *testing.T) {
	t.Setenv("BOOL_ENV", "true")

	result := getEnvBooleanOrDefault("BOOL_ENV", false)
	assert.Equal(t, result, true, "Result should be true")

	result = getEnvBooleanOrDefault("MISSING_BOOL_ENV", false)
	assert.Equal(t, result, false, "Result should be false")
}
