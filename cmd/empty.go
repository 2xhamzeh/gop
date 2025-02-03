package cmd

import (
	"github.com/2xhamzeh/gop/internal/template"
	"github.com/spf13/cobra"
)

var emptyCmd = &cobra.Command{
	Use:   "empty [dir-] [module-name]",
	Short: "setup empty project",
	Long:  `Setup an empty project with a predefined file structure.`,
	Example: `gop empty . github.com/user/project     # Create in current directory
gop empty myapp github.com/user/myapp   # Create in new directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return template.CreateFromTemplate(template.EmptyTemplate)
	},
}

func init() {
	rootCmd.AddCommand(emptyCmd)
}
