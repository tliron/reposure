package main

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	directclient "github.com/tliron/reposure/client/direct"
)

var logTo string
var verbose int
var colorize string
var host string
var certificatePath string
var username string
var password string
var token string

var roundTripper http.RoundTripper

func init() {
	rootCommand.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCommand.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCommand.PersistentFlags().StringVarP(&colorize, "colorize", "z", "true", "colorize output (boolean or \"force\")")
	rootCommand.PersistentFlags().StringVarP(&host, "host", "s", "localhost:5000", "registry host")
	rootCommand.PersistentFlags().StringVarP(&certificatePath, "certificate", "c", "", "registry TLS certificate file path (in PEM format)")
	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "", "registry authentication username")
	rootCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "registry authentication password")
	rootCommand.PersistentFlags().StringVarP(&token, "token", "t", "", "registry authentication token")
}

var rootCommand = &cobra.Command{
	Use:   "reposure-registry-client",
	Short: "Access a container image registry",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cleanup, err := terminal.ProcessColorizeFlag(colorize)
		util.FailOnError(err)
		if cleanup != nil {
			util.OnExitError(cleanup)
		}
		if logTo == "" {
			logging.Configure(verbose, nil)
		} else {
			logging.Configure(verbose, &logTo)
		}

		if host == "" {
			util.Fail("must provide \"--registry\"")
		}

		if certificatePath != "" {
			roundTripper, err = directclient.TLSRoundTripper(certificatePath)
			util.FailOnError(err)
		}
	},
}
