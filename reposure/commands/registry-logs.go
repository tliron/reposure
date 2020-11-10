package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	registryCommand.AddCommand(registryLogsCommand)
	registryLogsCommand.Flags().IntVarP(&tail, "tail", "t", -1, "number of most recent lines to print (<0 means all lines)")
	registryLogsCommand.Flags().BoolVarP(&follow, "follow", "f", false, "keep printing incoming logs")
}

var registryLogsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Show the logs of the Reposure container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		Logs("registry", "registry")
	},
}
