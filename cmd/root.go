/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gop",
	Short: "A tool for creating GO projects",
	Long: `A CLI tool for creating GO projects with a predefined structure.
You can create a new project with a single command and start working on your project immediately.
You can either create a simple empty project, a simple REST API project or an advanced REST API project`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gop.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
