package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
	registryclient "github.com/tliron/reposure/client/registry"
)

func init() {
	rootCommand.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List images in a container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		List(registry)
	},
}

func List(registry string) {
	images, err := registryclient.NewClient(roundTripper, username, password, token).ListImages(registry)
	util.FailOnError(err)
	for _, image := range images {
		fmt.Println(image)
	}
}
