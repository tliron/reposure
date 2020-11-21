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
	repositoryCommand.AddCommand(repositoryDescribeCommand)
}

var repositoryDescribeCommand = &cobra.Command{
	Use:   "describe [REPOSITORY NAME]",
	Short: "Describe a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		DescribeRepository(args[0])
	},
}

func DescribeRepository(repositoryName string) {
	// TODO: in cluster mode we must specify the namespace
	namespace := ""

	repository, err := NewClient().AdminClient().GetRepository(namespace, repositoryName)
	util.FailOnError(err)

	if format != "" {
		formatpkg.Print(resources.RepositoryToARD(repository), format, terminal.Stdout, strict, pretty)
	} else {
		fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("Name"), terminal.ColorValue(repository.Name))

		if repository.Spec.Direct != nil {
			fmt.Fprintf(terminal.Stdout, "  %s:\n", terminal.ColorTypeName("Direct"))
			if repository.Spec.Direct.Host != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.ColorTypeName("Host"), terminal.ColorValue(repository.Spec.Direct.Host))
			}
		}

		if repository.Spec.Indirect != nil {
			fmt.Fprintf(terminal.Stdout, "  %s:\n", terminal.ColorTypeName("Indirect"))
			if repository.Spec.Indirect.Namespace != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.ColorTypeName("Namespace"), terminal.ColorValue(repository.Spec.Indirect.Namespace))
			}
			if repository.Spec.Indirect.Service != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.ColorTypeName("Service"), terminal.ColorValue(repository.Spec.Indirect.Service))
			}
			fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.ColorTypeName("Port"), terminal.ColorValue(fmt.Sprintf("%d", repository.Spec.Indirect.Port)))
		}

		if repository.Spec.TLSSecret != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("TLSSecret"), terminal.ColorValue(repository.Spec.TLSSecret))
		}
		if repository.Spec.TLSSecretDataKey != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("TLSSecretDataKey"), terminal.ColorValue(repository.Spec.TLSSecretDataKey))
		}
		if repository.Spec.AuthSecret != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("AuthSecret"), terminal.ColorValue(repository.Spec.AuthSecret))
		}

		if repository.Status.SpoolerPod != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.ColorTypeName("SpoolerPod"), terminal.ColorValue(repository.Status.SpoolerPod))
		}
	}
}
