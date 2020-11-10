package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	imageCommand.AddCommand(imageListCommand)
}

var imageListCommand = &cobra.Command{
	Use:   "list [REPOSITORY NAME]",
	Short: "List images in a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repositoryName := args[0]

		adminClient := NewClient().AdminClient()
		repository, err := adminClient.GetRepository(namespace, repositoryName)
		util.FailOnError(err)
		commandClient, err := adminClient.CommandClient(repository)
		util.FailOnError(err)
		images, err := commandClient.ListImages()
		util.FailOnError(err)

		for _, image := range images {
			fmt.Fprintln(terminal.Stdout, image)
		}
	},
}
