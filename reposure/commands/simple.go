package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(simpleCommand)
}

var simpleCommand = &cobra.Command{
	Use:   "simple",
	Short: "Control the simple Reposure container image registry",
}
