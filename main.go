package main

import (
	"github.com/dmarro89/dare-db/server"
)

func main() {
	server := server.NewServerFactory().NewServer()
	server.Start()
	server.Stop()
}
