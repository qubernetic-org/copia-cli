package repo

import (
	"github.com/spf13/cobra"
	cloneCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/repo/clone"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/repo/list"
	viewCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/repo/view"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

// NewCmdRepo creates the `copia repo` command group.
func NewCmdRepo(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo <command>",
		Short: "Manage repositories",
		Long:  "Work with Copia repositories.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))
	cmd.AddCommand(viewCmd.NewCmdView(f))
	cmd.AddCommand(cloneCmd.NewCmdClone(f))

	return cmd
}
