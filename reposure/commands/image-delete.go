package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	imageCommand.AddCommand(imageDeleteCommand)
}

var imageDeleteCommand = &cobra.Command{
	Use:   "delete [REPOSITORY NAME] [IMAGE NAME]",
	Short: "Delete an image from a repository",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repositoryName := args[0]
		imageReference := args[1]

		adminClient := NewClient().AdminClient()
		repository, err := adminClient.GetRepository(namespace, repositoryName)
		util.FailOnError(err)
		spoolerClient := adminClient.SpoolerClient(repository)

		err = spoolerClient.DeleteImage(imageReference)
		util.FailOnError(err)
	},
}
