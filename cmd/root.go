/*
Copyright © 2025 2xhamzeh
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gop",
	Short: "CLI tool for creating GO projects",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print command usage") // Custom help flag
	rootCmd.PersistentFlags().MarkHidden("help")                               // Hide help flag
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})                       // Hide help command
	rootCmd.CompletionOptions.DisableDefaultCmd = true                         // Disable default completion command
}
