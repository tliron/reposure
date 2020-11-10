package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	registryCommand.AddCommand(registryShellCommand)
}

var registryShellCommand = &cobra.Command{
	Use:   "shell",
	Short: "Opens a shell to the Reposure container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		Shell("registry", "registry")
	},
}
