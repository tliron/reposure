package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	repositoryCommand.AddCommand(repositoryLogsCommand)
	repositoryLogsCommand.Flags().IntVarP(&tail, "tail", "t", -1, "number of most recent lines to print (<0 means all lines)")
	repositoryLogsCommand.Flags().BoolVarP(&follow, "follow", "f", false, "keep printing incoming logs")
}

var repositoryLogsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Show the logs of a repository surrogate",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		Logs("repository", "")
	},
}
