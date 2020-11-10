package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(repositoryCommand)
}

var repositoryCommand = &cobra.Command{
	Use:   "repository",
	Short: "Work with repositories",
}
