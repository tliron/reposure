package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var secure bool

func init() {
	registryCommand.AddCommand(registryInstallCommand)
	registryInstallCommand.Flags().BoolVarP(&cluster, "cluster", "c", false, "cluster mode")
	registryInstallCommand.Flags().StringVarP(&registry, "registry", "g", "docker.io", "registry address (use special value \"internal\" to discover internally deployed registry)")
	registryInstallCommand.Flags().BoolVarP(&secure, "secure", "s", true, "secure the registry (requires cert-manager)")
	registryInstallCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for installation to succeed")
}

var registryInstallCommand = &cobra.Command{
	Use:   "install",
	Short: "Install the Reposure container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		err := NewClient().AdminClient().InstallRegistry(registry, secure, wait)
		util.FailOnError(err)
	},
}
