package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tidyup",
	Short: "TidyUp is a CLI tool to clean up stale dev dependencies",
	Long:  `A fast and smart CLI tool to find and remove massive, unused folders like node_modules or target across your machine.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to TidyUp! Use 'tidyup --help' to see available commands.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}