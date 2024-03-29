package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	imageCommand.AddCommand(imageListCommand)
}

var imageListCommand = &cobra.Command{
	Use:   "list [REGISTRY NAME]",
	Short: "List images in a registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ListImages(args[0])
	},
}

func ListImages(registryName string) {
	adminClient := NewClient().AdminClient()
	registry, err := adminClient.GetRegistry(namespace, registryName)
	util.FailOnError(err)
	surrogateCommandClient, err := adminClient.SurrogateCommandClient(registry)
	util.FailOnError(err)
	imageNames, err := surrogateCommandClient.ListImages()
	util.FailOnError(err)

	for _, imageName := range imageNames {
		terminal.Println(imageName)
	}
}
