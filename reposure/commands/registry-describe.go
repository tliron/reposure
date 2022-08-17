package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
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
		transcribe.Print(resources.RegistryToARD(registry), format, os.Stdout, strict, pretty)
	} else {
		terminal.Printf("%s: %s\n", terminal.DefaultStylist.TypeName("Name"), terminal.DefaultStylist.Value(registry.Name))
		terminal.Printf("%s: %s\n", terminal.DefaultStylist.TypeName("Type"), terminal.DefaultStylist.Value(string(registry.Spec.Type)))

		if registry.Spec.Direct != nil {
			terminal.Printf("  %s:\n", terminal.DefaultStylist.TypeName("Direct"))
			if registry.Spec.Direct.Host != "" {
				terminal.Printf("    %s: %s\n", terminal.DefaultStylist.TypeName("Host"), terminal.DefaultStylist.Value(registry.Spec.Direct.Host))
			}
		}

		if registry.Spec.Indirect != nil {
			terminal.Printf("  %s:\n", terminal.DefaultStylist.TypeName("Indirect"))
			if registry.Spec.Indirect.Namespace != "" {
				terminal.Printf("    %s: %s\n", terminal.DefaultStylist.TypeName("Namespace"), terminal.DefaultStylist.Value(registry.Spec.Indirect.Namespace))
			}
			if registry.Spec.Indirect.Service != "" {
				terminal.Printf("    %s: %s\n", terminal.DefaultStylist.TypeName("Service"), terminal.DefaultStylist.Value(registry.Spec.Indirect.Service))
			}
			terminal.Printf("    %s: %s\n", terminal.DefaultStylist.TypeName("Port"), terminal.DefaultStylist.Value(fmt.Sprintf("%d", registry.Spec.Indirect.Port)))
		}

		if registry.Spec.AuthenticationSecret != "" {
			terminal.Printf("%s: %s\n", terminal.DefaultStylist.TypeName("AuthenticationSecret"), terminal.DefaultStylist.Value(registry.Spec.AuthenticationSecret))
		}
		if registry.Spec.AuthenticationSecretDataKey != "" {
			terminal.Printf("%s: %s\n", terminal.DefaultStylist.TypeName("AuthenticationSecretDataKey"), terminal.DefaultStylist.Value(registry.Spec.AuthenticationSecretDataKey))
		}
		if registry.Spec.AuthorizationSecret != "" {
			terminal.Printf("%s: %s\n", terminal.DefaultStylist.TypeName("AuthorizationSecret"), terminal.DefaultStylist.Value(registry.Spec.AuthorizationSecret))
		}

		if registry.Status.SurrogatePod != "" {
			terminal.Printf("%s: %s\n", terminal.DefaultStylist.TypeName("SurrogatePod"), terminal.DefaultStylist.Value(registry.Status.SurrogatePod))
		}
	}
}
