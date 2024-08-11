package main

import (
	"github.com/dmarro89/dare-db/auth"
	"github.com/dmarro89/dare-db/database"
	"github.com/dmarro89/dare-db/logger"
	"github.com/dmarro89/dare-db/server"
)

func main() {
	logger := logger.NewDareLogger()
	configuration := server.NewConfiguration("")
	database := database.NewDatabase()
	userStore := auth.NewUserStore()
	userStore.AddUser("admin", "password")
	dareServer := server.NewDareServer(database, userStore)
	server := server.NewFactory(configuration, logger).GetWebServer(dareServer)

	server.Start()
	server.Stop()
}
