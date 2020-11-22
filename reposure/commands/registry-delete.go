package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	registryCommand.AddCommand(registryDeleteCommand)
	registryDeleteCommand.Flags().BoolVarP(&all, "all", "a", false, "delete all registries")
}

var registryDeleteCommand = &cobra.Command{
	Use:   "delete [[REGISTRY NAME]]",
	Short: "Delete a registry",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			DeleteRegistry(args[0])
		} else if all {
			DeleteAllRegistries()
		} else {
			util.Fail("must provide registry name or specify \"--all\"")
		}
	},
}

func DeleteRegistry(registryName string) {
	// TODO: in cluster mode we must specify the namespace
	namespace := ""

	err := NewClient().AdminClient().DeleteRegistry(namespace, registryName)
	util.FailOnError(err)
}

func DeleteAllRegistries() {
	reposure := NewClient().AdminClient()
	registries, err := reposure.ListRegistries()
	util.FailOnError(err)
	if len(registries.Items) > 0 {
		for _, registry := range registries.Items {
			log.Infof("deleting registry: %s/%s", registry.Namespace, registry.Name)
			err := reposure.DeleteRegistry(registry.Namespace, registry.Name)
			util.FailOnError(err)
		}
	} else {
		log.Info("no registries to delete")
	}
}
