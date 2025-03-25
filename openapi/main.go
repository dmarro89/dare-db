package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	openapidocs "github.com/kohkimakimoto/echo-openapidocs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var workingDir string

const (
	URL_OPENAPI_DOC          = "/docs"
	URL_OPENAPI_FILE         = "/openapi-file"
	FILE_OPENAPI_DOC         = "daredb-openapi_primary.yaml"
	FILE_OPENAPI_DOC_READY   = "daredb-openapi_generated_webui.yaml"
	TEMPLATE_API_SERVER_NAME = "API_SERVER_NAME"
)

const (
	ECHO_LOG_FORMATTING = `${time_rfc3339}  remote_ip=${remote_ip}, method=${method}, status=${status}, uri=${uri}, host=${host}, error=${error}`
)

type customLogWriter struct {
}

type ConfigDatabase struct {
	OpenAPIServerHost       string `env:"OPENAPI_SERVER_HOST" env-default:"127.0.0.1"`
	OpenAPIServerPort       string `env:"OPENAPI_SERVER_PORT" env-default:"5002"`
	AccessOpenAPIServerHost string `env:"ACCESS_OPENAPI_SERVER_HOST" env-default:"127.0.0.1"`
	DareDBServerHost        string `env:"DAREDB_SERVER_HOST" env-default:"127.0.0.1"`
	DareDBServerPort        string `env:"DAREDB_SERVER_PORT" env-default:"5001"`
	DareDBWithTLS           bool   `env:"DAREDB_SERVER_WITH_TLS" env-default:"true"`
}

var Config ConfigDatabase

func serveOpenAPIDocYAML(c echo.Context) error {

	filePath := filepath.Join(workingDir, FILE_OPENAPI_DOC)
	filePathGen := filepath.Join(workingDir, FILE_OPENAPI_DOC_READY)
	dareDBServerName := fmt.Sprintf("http://%s:%s", Config.DareDBServerHost, Config.DareDBServerPort)

	if Config.DareDBWithTLS == true {
		dareDBServerName = fmt.Sprintf("https://%s:%s", Config.DareDBServerHost, Config.DareDBServerPort)
	}

	err := replaceInFile(filePath, filePathGen, TEMPLATE_API_SERVER_NAME, dareDBServerName)
	if err != nil {
		c.Echo().Logger.Error(err)
	}

	return c.File(filePathGen)
}

func GetEnvVariable(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Makes replacement of string in a file and creates a new file
func replaceInFile(infile, outfile, oldString, newString string) error {
	data, err := os.ReadFile(infile)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	newContent := strings.ReplaceAll(string(data), oldString, newString)
	return os.WriteFile(outfile, []byte(newContent), 0644)
}

// Redirects to the home page of the service
func serverNotFoundMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() == "/" {
			return c.Redirect(http.StatusMovedPermanently, URL_OPENAPI_DOC)
		}
		return next(c)
	}
}

// Making CORS less restrictive
func nonRestrictiveCORSMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//c.Response().Header().Set("Referrer-Policy", "origin-when-cross-origin") // Choose your desired policy
		c.Response().Header().Set("Referrer-Policy", "unsafe-url") // non secure workaround CORS
		return next(c)
	}
}

// Configures global value of the current working directory
func configWorkingDir() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting working directory: %v", err)
		return
	}
	log.Printf("Working directory: %v", workingDir)
}

func (writer customLogWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02 15:04:05") + " " + string(bytes))
}

func main() {

	log.SetFlags(0)
	log.SetOutput(new(customLogWriter))

	configWorkingDir()

	if err := cleanenv.ReadEnv(&Config); err != nil {
		log.Panicln("Was not able to read env")
	}

	//openapiServer := fmt.Sprintf("http://%s:%s", Config.OpenAPIServerHost, Config.OpenAPIServerPort)
	openapiFileServer := fmt.Sprintf("http://%s:%s", Config.AccessOpenAPIServerHost, Config.OpenAPIServerPort)
	//dareDBServer := fmt.Sprintf("http://%s:%s", Config.DareDBServerHost, Config.DareDBServerPort)
	//dareDBServerWithTLS := fmt.Sprintf("https://%s:%s", Config.DareDBServerHost, Config.DareDBServerPort)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(serverNotFoundMiddleware)
	e.Use(nonRestrictiveCORSMiddleware)

	// Add custom logging middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: ECHO_LOG_FORMATTING + "\n",
	}))

	// Enable CORS with specific origins (using v4 functions)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		//AllowOrigins: []string{openapiServer, openapiFileServer, dareDBServer, dareDBServerWithTLS},
	}))

	openapiServerFileUrl := openapiFileServer + URL_OPENAPI_FILE

	// Register the  Spotlight/Elements documentation with OpenAPI Spec url
	openapidocs.SwaggerUIDocuments(e, URL_OPENAPI_DOC, openapidocs.SwaggerUIConfig{
		SpecUrl: openapiServerFileUrl,
		Title:   "REST API",
	})

	// Serving OpenAPI 3.0 file
	e.GET(URL_OPENAPI_FILE, serveOpenAPIDocYAML)

	serverStartOn := fmt.Sprintf("%s:%s", Config.OpenAPIServerHost, Config.OpenAPIServerPort)
	log.Printf("Run server on: %s\n", serverStartOn)

	// Starting server
	e.Start(serverStartOn)
}
