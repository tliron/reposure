package commands

import (
	contextpkg "context"

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
		PushImage(contextpkg.TODO(), args[0], args[1], args[2])
	},
}

func PushImage(context contextpkg.Context, registryName string, imageName string, imagePath string) {
	adminClient := NewClient().AdminClient()
	registry, err := adminClient.GetRegistry(namespace, registryName)
	util.FailOnError(err)
	surrogateSpoolerClient := adminClient.SurrogateSpoolerClient(registry)

	// TODO:
	// 1) block until spooler picks up file
	// 2) forward errors from spooler

	err = surrogateSpoolerClient.PushTarballFromFile(context, imageName, imagePath)
	util.FailOnError(err)
}
