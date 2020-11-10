package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	registryCommand.AddCommand(registryUninstallCommand)
	registryUninstallCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for uninstallation to succeed")
}

var registryUninstallCommand = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the Reposure container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		NewClient().AdminClient().UninstallRegistry(wait)
	},
}
