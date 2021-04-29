package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	imageCommand.AddCommand(imageDeleteCommand)
	imageDeleteCommand.Flags().BoolVarP(&all, "all", "a", false, "delete all registries")
}

var imageDeleteCommand = &cobra.Command{
	Use:   "delete [REGISTRY NAME] [[IMAGE NAME]]",
	Short: "Delete images from a registry",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			DeleteImage(args[0], args[1])
		} else if all {
			DeleteAllImages(args[0])
		} else {
			util.Fail("must provide image name or specify \"--all\"")
		}
	},
}

func DeleteImage(registryName string, imageName string) {
	adminClient := NewClient().AdminClient()
	registry, err := adminClient.GetRegistry(namespace, registryName)
	util.FailOnError(err)
	surrogateSpoolerClient := adminClient.SurrogateSpoolerClient(registry)

	err = surrogateSpoolerClient.DeleteImage(imageName)
	util.FailOnError(err)
}

func DeleteAllImages(registryName string) {
	adminClient := NewClient().AdminClient()
	registry, err := adminClient.GetRegistry(namespace, registryName)
	util.FailOnError(err)
	surrogateCommandClient, err := adminClient.SurrogateCommandClient(registry)
	util.FailOnError(err)
	surrogateSpoolerClient := adminClient.SurrogateSpoolerClient(registry)

	imageNames, err := surrogateCommandClient.ListImages()
	util.FailOnError(err)
	if len(imageNames) > 0 {
		for _, imageName := range imageNames {
			log.Infof("deleting image: %s", imageName)
			err := surrogateSpoolerClient.DeleteImage(imageName)
			util.FailOnError(err)
		}
	} else {
		log.Info("no images to delete")
	}
}
