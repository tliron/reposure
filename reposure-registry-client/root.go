package main

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	registryclient "github.com/tliron/reposure/client/registry"
)

var logTo string
var verbose int
var colorize string
var registry string
var certificatePath string
var username string
var password string
var token string

var roundTripper http.RoundTripper

func init() {
	rootCommand.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCommand.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCommand.PersistentFlags().StringVarP(&colorize, "colorize", "z", "true", "colorize output (boolean or \"force\")")
	rootCommand.PersistentFlags().StringVarP(&registry, "registry", "r", "localhost:5000", "registry URL")
	rootCommand.PersistentFlags().StringVarP(&certificatePath, "certificate", "c", "", "registry TLS certificate file path (in PEM format)")
	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "", "registry authentication username")
	rootCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "registry authentication password")
	rootCommand.PersistentFlags().StringVarP(&token, "token", "t", "", "registry authentication token")
}

var rootCommand = &cobra.Command{
	Use:   "reposure-registry-client",
	Short: "Access a container image registry",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := terminal.ProcessColorizeFlag(colorize)
		util.FailOnError(err)

		if logTo == "" {
			util.ConfigureLogging(verbose, nil)
		} else {
			util.ConfigureLogging(verbose, &logTo)
		}

		if registry == "" {
			util.Fail("must provide \"--registry\"")
		}

		if certificatePath != "" {
			roundTripper, err = registryclient.TLSRoundTripper(certificatePath)
			util.FailOnError(err)
		}
	},
}
