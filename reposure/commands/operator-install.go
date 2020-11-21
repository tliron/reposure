package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	operatorCommand.AddCommand(operatorInstallCommand)
	operatorInstallCommand.Flags().BoolVarP(&cluster, "cluster", "c", false, "cluster mode")
	operatorInstallCommand.Flags().StringVarP(&registry, "registry", "g", "docker.io", "registry address (use special value \"internal\" to discover internally deployed registry)")
	operatorInstallCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for installation to succeed")
}

var operatorInstallCommand = &cobra.Command{
	Use:   "install",
	Short: "Install the Reposure operator",
	Run: func(cmd *cobra.Command, args []string) {
		err := NewClient().AdminClient().InstallOperator(registry, wait)
		util.FailOnError(err)
	},
}
