package main

import (
	"net/http"

	"github.com/dmarro89/dare-db/auth"
	"github.com/dmarro89/dare-db/database"
	"github.com/dmarro89/dare-db/logger"
	"github.com/dmarro89/dare-db/server"
	"github.com/rs/cors"
)

// ApplyCORS applies CORS to the provided handler.
func ApplyCORS(handler http.Handler, allowedOrigins []string) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})
	return c.Handler(handler)
}

func main() {
	logger := logger.NewDareLogger()
	configuration := server.NewConfiguration("")
	database := database.NewDatabase()
	userStore := auth.NewUserStore()
	userStore.AddUser(configuration.GetString("server.admin_user"), configuration.GetString("server.admin_password"))
	dareServer := server.NewDareServer(database, userStore)
	server := server.NewFactory(configuration, logger).GetWebServer(dareServer)

	server.Start()
	defer server.Stop()
}
