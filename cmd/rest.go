/*
Copyright © 2025 2xhamzeh
*/
package cmd

import (
	"errors"

	"github.com/2xhamzeh/gop/internal/template"
	"github.com/spf13/cobra"
)

var restCmd = &cobra.Command{
	Use:          "rest [module-name]",
	Short:        "Setup a REST API project",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing module name")
		} else if len(args) > 1 {
			return errors.New("too many arguments")
		}
		return template.CreateFromTemplate(template.RestTemplate, args[0])
	},
}

func init() {
	rootCmd.AddCommand(restCmd)
}
