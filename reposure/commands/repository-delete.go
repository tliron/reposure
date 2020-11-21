package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	repositoryCommand.AddCommand(repositoryDeleteCommand)
	repositoryDeleteCommand.Flags().BoolVarP(&all, "all", "a", false, "delete all repositories")
}

var repositoryDeleteCommand = &cobra.Command{
	Use:   "delete [[REPOSITORY NAME]]",
	Short: "Delete a repository",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			repositoryName := args[0]
			DeleteRepository(repositoryName)
		} else if all {
			DeleteAllRepositories()
		} else {
			util.Fail("must provide repository name or specify \"--all\"")
		}
	},
}

func DeleteRepository(repositoryName string) {
	// TODO: in cluster mode we must specify the namespace
	namespace := ""

	err := NewClient().AdminClient().DeleteRepository(namespace, repositoryName)
	util.FailOnError(err)
}

func DeleteAllRepositories() {
	reposure := NewClient().AdminClient()
	repositories, err := reposure.ListRepositories()
	util.FailOnError(err)
	for _, repository := range repositories.Items {
		log.Infof("deleting repository: %s/%s", repository.Namespace, repository.Name)
		err := reposure.DeleteRepository(repository.Namespace, repository.Name)
		util.FailOnError(err)
	}
}
