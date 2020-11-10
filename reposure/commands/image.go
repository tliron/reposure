package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(imageCommand)
}

var imageCommand = &cobra.Command{
	Use:   "image",
	Short: "Work with images",
}
