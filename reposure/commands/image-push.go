package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	imageCommand.AddCommand(imagePushCommand)
}

var imagePushCommand = &cobra.Command{
	Use:   "push [REPOSITORY NAME] [IMAGE NAME] [IMAGE PATH]",
	Short: "Push an image to a repository",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		repositoryName := args[0]
		imageReference := args[1]
		imagePath := args[2]

		adminClient := NewClient().AdminClient()
		repository, err := adminClient.GetRepository(namespace, repositoryName)
		util.FailOnError(err)
		spoolerClient := adminClient.SpoolerClient(repository)

		file, err := os.Open(imagePath)
		util.FailOnError(err)
		defer file.Close()
		err = spoolerClient.PushImage(imageReference, file)
		util.FailOnError(err)
	},
}
