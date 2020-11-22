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
		fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("Name"), terminal.ColorValue(registry.Name))

		if registry.Spec.Direct != nil {
			fmt.Fprintf(terminal.Stdout, "  %s:\n", terminal.ColorTypeName("Direct"))
			if registry.Spec.Direct.Host != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.ColorTypeName("Host"), terminal.ColorValue(registry.Spec.Direct.Host))
			}
		}

		if registry.Spec.Indirect != nil {
			fmt.Fprintf(terminal.Stdout, "  %s:\n", terminal.ColorTypeName("Indirect"))
			if registry.Spec.Indirect.Namespace != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.ColorTypeName("Namespace"), terminal.ColorValue(registry.Spec.Indirect.Namespace))
			}
			if registry.Spec.Indirect.Service != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.ColorTypeName("Service"), terminal.ColorValue(registry.Spec.Indirect.Service))
			}
			fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.ColorTypeName("Port"), terminal.ColorValue(fmt.Sprintf("%d", registry.Spec.Indirect.Port)))
		}

		if registry.Spec.AuthenticationSecret != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("AuthenticationSecret"), terminal.ColorValue(registry.Spec.AuthenticationSecret))
		}
		if registry.Spec.AuthenticationSecretDataKey != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("AuthenticationSecretDataKey"), terminal.ColorValue(registry.Spec.AuthenticationSecretDataKey))
		}
		if registry.Spec.AuthorizationSecret != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("AuthorizationSecret"), terminal.ColorValue(registry.Spec.AuthorizationSecret))
		}

		if registry.Status.SurrogatePod != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("SurrogatePod"), terminal.ColorValue(registry.Status.SurrogatePod))
		}
	}
}
