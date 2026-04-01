package issue

import (
	"github.com/spf13/cobra"
	closeCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue/close"
	commentCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue/comment"
	createCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue/create"
	editCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue/edit"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue/list"
	viewCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue/view"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

// NewCmdIssue creates the `copia issue` command group.
func NewCmdIssue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue <command>",
		Short: "Manage issues",
		Long:  "Work with Copia repository issues.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))
	cmd.AddCommand(createCmd.NewCmdCreate(f))
	cmd.AddCommand(viewCmd.NewCmdView(f))
	cmd.AddCommand(closeCmd.NewCmdClose(f))
	cmd.AddCommand(commentCmd.NewCmdComment(f))
	cmd.AddCommand(editCmd.NewCmdEdit(f))

	return cmd
}
