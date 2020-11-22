package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	simpleCommand.AddCommand(simpleShellCommand)
}

var simpleShellCommand = &cobra.Command{
	Use:   "shell",
	Short: "Opens a shell to the simple Reposure container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		Shell("simple", "registry")
	},
}
