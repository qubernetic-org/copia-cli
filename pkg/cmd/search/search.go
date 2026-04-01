package search

import (
	"github.com/spf13/cobra"
	issuesCmd "github.com/qubernetic/copia-cli/pkg/cmd/search/issues"
	reposCmd "github.com/qubernetic/copia-cli/pkg/cmd/search/repos"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
)

func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <command>",
		Short: "Search across Copia",
		Long:  "Search repositories and issues.",
	}

	cmd.AddCommand(reposCmd.NewCmdSearchRepos(f))
	cmd.AddCommand(issuesCmd.NewCmdSearchIssues(f))

	return cmd
}
