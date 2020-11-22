package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	simpleCommand.AddCommand(simpleLogsCommand)
	simpleLogsCommand.Flags().IntVarP(&tail, "tail", "t", -1, "number of most recent lines to print (<0 means all lines)")
	simpleLogsCommand.Flags().BoolVarP(&follow, "follow", "f", false, "keep printing incoming logs")
}

var simpleLogsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Show the logs of the simple Reposure container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		Logs("simple", "registry")
	},
}
