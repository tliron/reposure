package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var secure bool

func init() {
	simpleCommand.AddCommand(simpleInstallCommand)
	simpleInstallCommand.Flags().BoolVarP(&clusterMode, "cluster", "c", false, "cluster mode")
	simpleInstallCommand.Flags().StringVarP(&sourceRegistry, "registry", "g", "docker.io", "source registry host (use special value \"internal\" to discover internally deployed simple)")
	simpleInstallCommand.Flags().BoolVarP(&secure, "secure", "s", true, "secure the registry (requires cert-manager)")
	simpleInstallCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for installation to succeed")
}

var simpleInstallCommand = &cobra.Command{
	Use:   "install",
	Short: "Install the simple Reposure container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		InstallRegistry()
	},
}

func InstallRegistry() {
	err := NewClient().AdminClient().InstallRegistry(sourceRegistry, secure, wait)
	util.FailOnError(err)
}
