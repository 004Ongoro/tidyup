package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tidyup",
	Short: "TidyUp: A smart manager for your developer 'junk' folders",
	Long: `TidyUp is a high-performance CLI tool designed to reclaim disk space. 
It identifies stale dependencies (like node_modules, target, .venv) 
based on the last time you actually worked on the project.

Safe by default: It automatically skips system folders and IDE configurations.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = "0.1.0"
}