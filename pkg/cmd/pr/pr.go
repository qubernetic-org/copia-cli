package pr

import (
	"github.com/spf13/cobra"
	checkoutCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/checkout"
	closeCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/close"
	createCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/create"
	diffCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/diff"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/list"
	mergeCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/merge"
	reviewCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/review"
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
	cmd.AddCommand(reviewCmd.NewCmdReview(f))
	cmd.AddCommand(diffCmd.NewCmdDiff(f))
	cmd.AddCommand(checkoutCmd.NewCmdCheckout(f))

	return cmd
}
