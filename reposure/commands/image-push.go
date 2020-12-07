package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	imageCommand.AddCommand(imagePushCommand)
}

var imagePushCommand = &cobra.Command{
	Use:   "push [REGISTRY NAME] [IMAGE NAME] [IMAGE PATH]",
	Short: "Push an image to a registry",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		PushImage(args[0], args[1], args[2])
	},
}

func PushImage(registryName string, imageName string, imagePath string) {
	adminClient := NewClient().AdminClient()
	registry, err := adminClient.GetRegistry(namespace, registryName)
	util.FailOnError(err)
	spoolerClient := adminClient.SpoolerClient(registry)

	err = spoolerClient.PushTarballFromFile(imageName, imagePath)
	util.FailOnError(err)
}
