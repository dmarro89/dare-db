package main

import (
	"github.com/dmarro89/dare-db/database"
	"github.com/dmarro89/dare-db/server"
	"github.com/dmarro89/dare-db/logger"
)


func main() {
	darelogger := darelog.NewLogger(darelog.GetEnvLOGLEVEL())

	database := database.NewDatabase(darelogger)
	database.Dict.GenerateRandomBytes()

	dareServer := server.NewDareServer(database, darelogger)
	server := server.NewFactory().GetWebServer(dareServer)

	server.Start()
	server.Stop()
}
