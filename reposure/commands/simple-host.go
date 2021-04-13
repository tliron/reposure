package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	simpleCommand.AddCommand(simpleHostCommand)
}

var simpleHostCommand = &cobra.Command{
	Use:   "host",
	Short: "Get the host of the simple Reposure container image registry",
	Run: func(cmd *cobra.Command, args []string) {
		host, err := NewClient().AdminClient().SimpleHost()
		util.FailOnError(err)
		terminal.Println(host)
	},
}
