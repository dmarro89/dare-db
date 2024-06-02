package main

import (
	"github.com/dmarro89/dare-db/database"
	"github.com/dmarro89/dare-db/server"
)

func main() {

	server.Configure("config.toml")

	database := database.NewDatabase()
	dareServer := server.NewDareServer(database)
	server := server.NewFactory().GetWebServer(dareServer)

	server.Start()
	server.Stop()
}
