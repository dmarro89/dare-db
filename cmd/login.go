package cmd

import (
	"errors"
	"fmt"

	"github.com/dmarro89/dare-db/auth"
	"github.com/spf13/cobra"
)

type CLIAutorhizer struct {
	UsersStore *auth.UserStore
}

var CLIAuthorizer = &CLIAutorhizer{
	UsersStore: auth.NewUserStore(),
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into the application",
	Run: func(cmd *cobra.Command, args []string) {
		username, password := args[0], args[1]
		err := CLIAuthorizer.login(username, password) // Funzione che gestisce l'autenticazione
		if err != nil {
			fmt.Println("Login failed:", err)
			return
		}

		fmt.Println("Login successful! Session token saved.")
	},
}

func (cliAuthorizer *CLIAutorhizer) login(username, password string) error {
	if cliAuthorizer.UsersStore.ValidateCredentials(username, password) {
		return nil

	}
	return errors.New("invalid credentials")
}
