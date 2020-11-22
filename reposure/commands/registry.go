package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(registryCommand)
}

var registryCommand = &cobra.Command{
	Use:   "registry",
	Short: "Work with registries",
}
