package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
)

func init() {
	registryCommand.AddCommand(registryDescribeCommand)
}

var registryDescribeCommand = &cobra.Command{
	Use:   "describe [REGISTRY NAME]",
	Short: "Describe a registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		DescribeRegistry(args[0])
	},
}

func DescribeRegistry(registryName string) {
	// TODO: in cluster mode we must specify the namespace
	namespace := ""

	registry, err := NewClient().AdminClient().GetRegistry(namespace, registryName)
	util.FailOnError(err)

	if format != "" {
		Transcriber().Write(resources.RegistryToARD(registry))
	} else {
		terminal.Printf("%s: %s\n", terminal.StdoutStylist.TypeName("Name"), terminal.StdoutStylist.Value(registry.Name))
		terminal.Printf("%s: %s\n", terminal.StdoutStylist.TypeName("Type"), terminal.StdoutStylist.Value(string(registry.Spec.Type)))

		if registry.Spec.Direct != nil {
			terminal.Printf("  %s:\n", terminal.StdoutStylist.TypeName("Direct"))
			if registry.Spec.Direct.Host != "" {
				terminal.Printf("    %s: %s\n", terminal.StdoutStylist.TypeName("Host"), terminal.StdoutStylist.Value(registry.Spec.Direct.Host))
			}
		}

		if registry.Spec.Indirect != nil {
			terminal.Printf("  %s:\n", terminal.StdoutStylist.TypeName("Indirect"))
			if registry.Spec.Indirect.Namespace != "" {
				terminal.Printf("    %s: %s\n", terminal.StdoutStylist.TypeName("Namespace"), terminal.StdoutStylist.Value(registry.Spec.Indirect.Namespace))
			}
			if registry.Spec.Indirect.Service != "" {
				terminal.Printf("    %s: %s\n", terminal.StdoutStylist.TypeName("Service"), terminal.StdoutStylist.Value(registry.Spec.Indirect.Service))
			}
			terminal.Printf("    %s: %s\n", terminal.StdoutStylist.TypeName("Port"), terminal.StdoutStylist.Value(fmt.Sprintf("%d", registry.Spec.Indirect.Port)))
		}

		if registry.Spec.AuthenticationSecret != "" {
			terminal.Printf("%s: %s\n", terminal.StdoutStylist.TypeName("AuthenticationSecret"), terminal.StdoutStylist.Value(registry.Spec.AuthenticationSecret))
		}
		if registry.Spec.AuthenticationSecretDataKey != "" {
			terminal.Printf("%s: %s\n", terminal.StdoutStylist.TypeName("AuthenticationSecretDataKey"), terminal.StdoutStylist.Value(registry.Spec.AuthenticationSecretDataKey))
		}
		if registry.Spec.AuthorizationSecret != "" {
			terminal.Printf("%s: %s\n", terminal.StdoutStylist.TypeName("AuthorizationSecret"), terminal.StdoutStylist.Value(registry.Spec.AuthorizationSecret))
		}

		if registry.Status.SurrogatePod != "" {
			terminal.Printf("%s: %s\n", terminal.StdoutStylist.TypeName("SurrogatePod"), terminal.StdoutStylist.Value(registry.Status.SurrogatePod))
		}
	}
}
