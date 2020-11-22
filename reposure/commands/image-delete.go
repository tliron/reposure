package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	imageCommand.AddCommand(imageDeleteCommand)
}

var imageDeleteCommand = &cobra.Command{
	Use:   "delete [REGISTRY NAME] [IMAGE NAME]",
	Short: "Delete an image from a registry",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		DeleteImage(args[0], args[1])
	},
}

func DeleteImage(registryName string, imageName string) {
	adminClient := NewClient().AdminClient()
	registry, err := adminClient.GetRegistry(namespace, registryName)
	util.FailOnError(err)
	spoolerClient := adminClient.SpoolerClient(registry)

	err = spoolerClient.DeleteImage(imageName)
	util.FailOnError(err)
}
