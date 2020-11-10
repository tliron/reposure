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
	Use:   "pull [REPOSITORY NAME] [IMAGE NAME]",
	Short: "Pull an image from a repository",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repositoryName := args[0]
		imageReference := args[1]

		adminClient := NewClient().AdminClient()
		repository, err := adminClient.GetRepository(namespace, repositoryName)
		util.FailOnError(err)
		commandClient, err := adminClient.CommandClient(repository)
		util.FailOnError(err)

		if unpack {
			err = commandClient.PullLayer(imageReference, os.Stdout)
		} else {
			err = commandClient.PullTarball(imageReference, os.Stdout)
		}
		util.FailOnError(err)
	},
}
