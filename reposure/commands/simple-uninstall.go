package commands

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	simpleCommand.AddCommand(simpleUninstallCommand)
	simpleUninstallCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for uninstallation to succeed")
}

var simpleUninstallCommand = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the simple Reposure container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		UninstallRegistry()
	},
}

func UninstallRegistry() {
	NewClient().AdminClient().UninstallSimple(context.TODO(), wait)
}
