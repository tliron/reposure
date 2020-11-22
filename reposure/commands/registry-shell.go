package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	registryCommand.AddCommand(registryShellCommand)
}

var registryShellCommand = &cobra.Command{
	Use:   "shell [REGISTRY NAME]",
	Short: "Opens a shell to a registry surrogate",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Shell(fmt.Sprintf("surrogate-%s", args[0]), "surrogate")
	},
}
