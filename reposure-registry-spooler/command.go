package main

import (
	"github.com/spf13/cobra"
	cobrautil "github.com/tliron/kutil/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

var logTo string
var verbose int
var colorize string
var directoryPath string
var registry string
var certificatePath string
var username string
var password string
var token string
var queue int
var healthPort uint

func init() {
	command.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	command.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	command.PersistentFlags().StringVarP(&colorize, "colorize", "z", "true", "colorize output (boolean or \"force\")")
	command.PersistentFlags().StringVarP(&directoryPath, "directory", "d", "/spool", "spool directory path")
	command.PersistentFlags().StringVarP(&registry, "registry", "r", "localhost:5000", "registry URL")
	command.PersistentFlags().StringVarP(&certificatePath, "certificate", "c", "", "registry TLS certificate file path (in PEM format)")
	command.PersistentFlags().StringVarP(&username, "username", "u", "", "registry authentication username")
	command.PersistentFlags().StringVarP(&password, "password", "p", "", "registry authentication password")
	command.PersistentFlags().StringVarP(&token, "token", "t", "", "registry authentication token")
	command.PersistentFlags().IntVarP(&queue, "queue", "q", 10, "maximum number of files to queue at once")
	command.PersistentFlags().UintVar(&healthPort, "health-port", 8086, "HTTP port for health check (for liveness and readiness probes)")

	cobrautil.SetFlagsFromEnvironment("REPOSURE_REGISTRY_SPOOLER_", command)
}

var command = &cobra.Command{
	Use:   "reposure-registry-spooler",
	Short: "Spooler for a container image registry",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := terminal.ProcessColorizeFlag(colorize)
		util.FailOnError(err)

		if logTo == "" {
			util.ConfigureLogging(verbose, nil)
		} else {
			util.ConfigureLogging(verbose, &logTo)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if (directoryPath == "") || (registry == "") {
			util.Fail("must provide \"--directory\" and \"--registry\"")
		}

		RunSpooler(registry, directoryPath)
	},
}
