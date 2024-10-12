/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dare-db",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var command1Cmd = &cobra.Command{
	Use:   "command1",
	Short: "Esegui command1",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Eseguito command1")
	},
}

// Definisci il comando command2
var command2Cmd = &cobra.Command{
	Use:   "command2",
	Short: "Esegui command2",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Eseguito command2")
	},
}

func init() {
	RootCmd.AddCommand(command1Cmd)
	RootCmd.AddCommand(command2Cmd)
	RootCmd.AddCommand(LoginCmd)
}
