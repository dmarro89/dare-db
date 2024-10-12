package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dmarro89/dare-db/auth"
	"github.com/dmarro89/dare-db/cmd"
	"github.com/dmarro89/dare-db/database"
	"github.com/dmarro89/dare-db/logger"
	"github.com/dmarro89/dare-db/server"
)

func main() {
	appMode := os.Getenv("APP_MODE")
	fmt.Println(`app mode:`, appMode)

	userStore := auth.NewUserStore()
	configuration := server.NewConfiguration("")
	userStore.AddUser(configuration.GetString("server.admin_user"), configuration.GetString("server.admin_password"))

	switch appMode {
	case "cli":
		runCLIPrompt()
	default:
		logger := logger.NewDareLogger()
		database := database.NewDatabase()
		dareServer := server.NewDareServer(database, userStore)
		server := server.NewFactory(configuration, logger).GetWebServer(dareServer)
		server.Start()
		defer server.Stop()
	}
}

func runCLIPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command: ")

		// Leggi il comando inserito dall'utente
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		// Se l'utente inserisce 'exit', esci dalla CLI
		if command == "exit" {
			fmt.Println("Exiting CLI...")
			break
		}

		// Esegui il comando con Cobra
		args := strings.Split(command, " ")
		cmd.RootCmd.SetArgs(args)

		if err := cmd.RootCmd.Execute(); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
