package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var unpack bool

func init() {
	imageCommand.AddCommand(imagePullCommand)
	imagePullCommand.Flags().BoolVarP(&unpack, "unpack", "u", false, "untar tarball and gunzip first layer")
}

var imagePullCommand = &cobra.Command{
	Use:   "pull [REGISTRY NAME] [IMAGE NAME]",
	Short: "Pull an image from a registry",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		PullImage(args[0], args[1])
	},
}

func PullImage(registryName string, imageName string) {
	adminClient := NewClient().AdminClient()
	registry, err := adminClient.GetRegistry(namespace, registryName)
	util.FailOnError(err)
	commandClient, err := adminClient.CommandClient(registry)
	util.FailOnError(err)

	if unpack {
		err = commandClient.PullLayer(imageName, os.Stdout)
	} else {
		err = commandClient.PullTarball(imageName, os.Stdout)
	}
	util.FailOnError(err)
}
