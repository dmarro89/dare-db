package server

import (
	"fmt"
	"os"
	"strings"
)

const DARE_PORT = "DARE_PORT"
const DARE_HOST = "DARE_HOST"
const DARE_TLS_ENABLED = "DARE_TLS_ENABLED"
const DARE_TLS_CERT_PRIVATE  = "DARE_TLS_CERT_PRIVATE"
const DARE_TLS_CERT_PUBLIC = "DARE_TLS_CERT_PUBLIC"

type Configuration struct {
	Port        string
	Host        string
	TLSEnabled  bool
	TLSCertFile string
	TLSKeyFile  string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:        fmt.Sprintf("%s", getEnvOrDefault(DARE_PORT, "2605")),
		Host:        fmt.Sprintf("%s", getEnvOrDefault(DARE_HOST, "")),
		TLSEnabled:  getEnvBooleanOrDefault(DARE_TLS_ENABLED, false),
		TLSCertFile: fmt.Sprintf("%s", getEnvOrDefault(DARE_TLS_CERT_PRIVATE , "")),
		TLSKeyFile:  fmt.Sprintf("%s", getEnvOrDefault(DARE_TLS_CERT_PUBLIC, "")),
	}
}

func getEnvOrDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBooleanOrDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.EqualFold(value, "true")
}
