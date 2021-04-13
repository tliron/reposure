package commands

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	registryCommand.AddCommand(registryInfoCommand)
}

var registryInfoCommand = &cobra.Command{
	Use:   "info [REGISTRY NAME] [host|username|password|token|cert]",
	Short: "Get registry information",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		RegistryInfo(args[0], args[1])
	},
}

func RegistryInfo(registryName string, field string) {
	// TODO: in cluster mode we must specify the namespace
	namespace := ""

	registryClient := NewClient().AdminClient().RegistryClient()
	registry, err := registryClient.Get(namespace, registryName)
	util.FailOnError(err)

	switch field {
	case "host":
		host, err := registryClient.GetHost(registry)
		util.FailOnError(err)
		if host != "" {
			terminal.Println(host)
		}

	case "username":
		_, username, _, _, err := registryClient.GetAuthorization(registry)
		util.FailOnError(err)
		if username != "" {
			terminal.Println(username)
		}

	case "password":
		_, _, password, _, err := registryClient.GetAuthorization(registry)
		util.FailOnError(err)
		if password != "" {
			terminal.Println(password)
		}

	case "token":
		_, _, _, token, err := registryClient.GetAuthorization(registry)
		util.FailOnError(err)
		if token != "" {
			terminal.Println(token)
		}

	case "cert":
		cert, err := registryClient.GetTLSCertBytes(registry)
		util.FailOnError(err)
		cert_ := strings.Trim(util.BytesToString(cert), "\n")
		if cert_ != "" {
			terminal.Println(cert_)
		}

	default:
		util.Failf("unsupported field: %s", field)
	}
}
