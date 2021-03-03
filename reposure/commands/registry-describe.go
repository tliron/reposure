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
		fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.StyleTypeName("Name"), terminal.StyleValue(registry.Name))
		fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.StyleTypeName("Type"), terminal.StyleValue(string(registry.Spec.Type)))

		if registry.Spec.Direct != nil {
			fmt.Fprintf(terminal.Stdout, "  %s:\n", terminal.StyleTypeName("Direct"))
			if registry.Spec.Direct.Host != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.StyleTypeName("Host"), terminal.StyleValue(registry.Spec.Direct.Host))
			}
		}

		if registry.Spec.Indirect != nil {
			fmt.Fprintf(terminal.Stdout, "  %s:\n", terminal.StyleTypeName("Indirect"))
			if registry.Spec.Indirect.Namespace != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.StyleTypeName("Namespace"), terminal.StyleValue(registry.Spec.Indirect.Namespace))
			}
			if registry.Spec.Indirect.Service != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.StyleTypeName("Service"), terminal.StyleValue(registry.Spec.Indirect.Service))
			}
			fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.StyleTypeName("Port"), terminal.StyleValue(fmt.Sprintf("%d", registry.Spec.Indirect.Port)))
		}

		if registry.Spec.AuthenticationSecret != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.StyleTypeName("AuthenticationSecret"), terminal.StyleValue(registry.Spec.AuthenticationSecret))
		}
		if registry.Spec.AuthenticationSecretDataKey != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.StyleTypeName("AuthenticationSecretDataKey"), terminal.StyleValue(registry.Spec.AuthenticationSecretDataKey))
		}
		if registry.Spec.AuthorizationSecret != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.StyleTypeName("AuthorizationSecret"), terminal.StyleValue(registry.Spec.AuthorizationSecret))
		}

		if registry.Status.SurrogatePod != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.StyleTypeName("SurrogatePod"), terminal.StyleValue(registry.Status.SurrogatePod))
		}
	}
}
