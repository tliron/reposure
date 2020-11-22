package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	registryCommand.AddCommand(registryLogsCommand)
	registryLogsCommand.Flags().IntVarP(&tail, "tail", "t", -1, "number of most recent lines to print (<0 means all lines)")
	registryLogsCommand.Flags().BoolVarP(&follow, "follow", "f", false, "keep printing incoming logs")
}

var registryLogsCommand = &cobra.Command{
	Use:   "logs [REGISTRY NAME]",
	Short: "Show the logs of a registry surrogate",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Logs(fmt.Sprintf("surrogate-%s", args[0]), "surrogate")
	},
}
