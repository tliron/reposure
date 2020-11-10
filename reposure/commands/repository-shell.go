package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	repositoryCommand.AddCommand(repositoryShellCommand)
}

var repositoryShellCommand = &cobra.Command{
	Use:   "shell",
	Short: "Opens a shell to a repository surrogate",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		Shell("repository", "")
	},
}
