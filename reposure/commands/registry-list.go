package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
	"github.com/tliron/kutil/util"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
)

func init() {
	registryCommand.AddCommand(registryListCommand)
}

var registryListCommand = &cobra.Command{
	Use:   "list",
	Short: "List registries",
	Run: func(cmd *cobra.Command, args []string) {
		ListRegistries()
	},
}

func ListRegistries() {
	registries, err := NewClient().AdminClient().ListRegistries()
	util.FailOnError(err)
	if len(registries.Items) == 0 {
		return
	}
	// TODO: sort registries by name? they seem already sorted!

	switch format {
	case "":
		table := terminal.NewTable(maxWidth, "Name", "Host", "Namespace", "Service", "Port", "SurrogatePod")
		for _, registry := range registries.Items {
			if registry.Spec.Direct != nil {
				table.Add(registry.Name, registry.Spec.Direct.Host, "", "", "", registry.Status.SurrogatePod)
			} else if registry.Spec.Indirect != nil {
				table.Add(registry.Name, "", registry.Spec.Indirect.Namespace, registry.Spec.Indirect.Service, fmt.Sprintf("%d", registry.Spec.Indirect.Port), registry.Status.SurrogatePod)
			}
		}
		table.Print()

	case "bare":
		for _, registry := range registries.Items {
			terminal.Println(registry.Name)
		}

	default:
		list := make(ard.List, len(registries.Items))
		for index, registry := range registries.Items {
			list[index] = resources.RegistryToARD(&registry)
		}
		transcribe.Print(list, format, os.Stdout, strict, pretty)
	}
}
