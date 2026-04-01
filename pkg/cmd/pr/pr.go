package pr

import (
	"github.com/spf13/cobra"
	closeCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/close"
	createCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/create"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/list"
	mergeCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/merge"
	viewCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/view"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

// NewCmdPR creates the `copia pr` command group.
func NewCmdPR(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr <command>",
		Short: "Manage pull requests",
		Long:  "Work with Copia pull requests.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))
	cmd.AddCommand(createCmd.NewCmdCreate(f))
	cmd.AddCommand(viewCmd.NewCmdView(f))
	cmd.AddCommand(mergeCmd.NewCmdMerge(f))
	cmd.AddCommand(closeCmd.NewCmdClose(f))

	return cmd
}
