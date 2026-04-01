package release

import (
	"github.com/spf13/cobra"
	createCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/release/create"
	deleteCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/release/delete"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/release/list"
	uploadCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/release/upload"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

func NewCmdRelease(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release <command>",
		Short: "Manage releases",
		Long:  "Work with Copia repository releases.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))
	cmd.AddCommand(createCmd.NewCmdCreate(f))
	cmd.AddCommand(deleteCmd.NewCmdDelete(f))
	cmd.AddCommand(uploadCmd.NewCmdUpload(f))

	return cmd
}
