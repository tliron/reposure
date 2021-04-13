package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	formatpkg "github.com/tliron/kutil/format"
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
		formatpkg.Print(resources.RegistryToARD(registry), format, terminal.Stdout, strict, pretty)
	} else {
		terminal.Printf("%s: %s\n", terminal.Stylize.TypeName("Name"), terminal.Stylize.Value(registry.Name))
		terminal.Printf("%s: %s\n", terminal.Stylize.TypeName("Type"), terminal.Stylize.Value(string(registry.Spec.Type)))

		if registry.Spec.Direct != nil {
			terminal.Printf("  %s:\n", terminal.Stylize.TypeName("Direct"))
			if registry.Spec.Direct.Host != "" {
				terminal.Printf("    %s: %s\n", terminal.Stylize.TypeName("Host"), terminal.Stylize.Value(registry.Spec.Direct.Host))
			}
		}

		if registry.Spec.Indirect != nil {
			terminal.Printf("  %s:\n", terminal.Stylize.TypeName("Indirect"))
			if registry.Spec.Indirect.Namespace != "" {
				terminal.Printf("    %s: %s\n", terminal.Stylize.TypeName("Namespace"), terminal.Stylize.Value(registry.Spec.Indirect.Namespace))
			}
			if registry.Spec.Indirect.Service != "" {
				terminal.Printf("    %s: %s\n", terminal.Stylize.TypeName("Service"), terminal.Stylize.Value(registry.Spec.Indirect.Service))
			}
			terminal.Printf("    %s: %s\n", terminal.Stylize.TypeName("Port"), terminal.Stylize.Value(fmt.Sprintf("%d", registry.Spec.Indirect.Port)))
		}

		if registry.Spec.AuthenticationSecret != "" {
			terminal.Printf("%s: %s\n", terminal.Stylize.TypeName("AuthenticationSecret"), terminal.Stylize.Value(registry.Spec.AuthenticationSecret))
		}
		if registry.Spec.AuthenticationSecretDataKey != "" {
			terminal.Printf("%s: %s\n", terminal.Stylize.TypeName("AuthenticationSecretDataKey"), terminal.Stylize.Value(registry.Spec.AuthenticationSecretDataKey))
		}
		if registry.Spec.AuthorizationSecret != "" {
			terminal.Printf("%s: %s\n", terminal.Stylize.TypeName("AuthorizationSecret"), terminal.Stylize.Value(registry.Spec.AuthorizationSecret))
		}

		if registry.Status.SurrogatePod != "" {
			terminal.Printf("%s: %s\n", terminal.Stylize.TypeName("SurrogatePod"), terminal.Stylize.Value(registry.Status.SurrogatePod))
		}
	}
}
