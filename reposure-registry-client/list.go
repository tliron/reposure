package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	rootCommand.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List images in a container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		List()
	},
}

func List() {
	images, err := NewClient().ListImages()
	util.FailOnError(err)
	for _, image := range images {
		fmt.Println(image)
	}
}
