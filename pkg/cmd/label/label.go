package label

import (
	"github.com/spf13/cobra"
	createCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/label/create"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/label/list"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

// NewCmdLabel creates the `copia label` command group.
func NewCmdLabel(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label <command>",
		Short: "Manage labels",
		Long:  "Work with repository labels.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))
	cmd.AddCommand(createCmd.NewCmdCreate(f))

	return cmd
}
