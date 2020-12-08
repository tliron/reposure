package main

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var output string

func init() {
	rootCommand.AddCommand(pullCommand)
	pullCommand.PersistentFlags().StringVarP(&output, "output", "o", "", "output to file (defaults to stdout)")
}

var pullCommand = &cobra.Command{
	Use:   "pull [IMAGE NAME]",
	Short: "Pull tarball from a container image registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Pull(args[0])
	},
}

func Pull(imageName string) {
	err := NewClient().PullTarball(imageName, output)
	util.FailOnError(err)
}
