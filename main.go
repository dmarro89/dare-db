package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dmarro89/dare-db/database"
	"github.com/dmarro89/dare-db/server"
)

func main() {
	srv := server.NewServer((database.NewDatabase()))

	http.HandleFunc("/get", srv.HandlerGet)
	http.HandleFunc("/set", srv.HandlerSet)
	http.HandleFunc("/delete", srv.HandlerDelete)

	port := getEnvOrDefault("DARE_PORT", 2605)

	fmt.Printf("Server listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}

}

func getEnvOrDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
