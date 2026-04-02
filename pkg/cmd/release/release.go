package release

import (
	"github.com/spf13/cobra"
	createCmd "github.com/qubernetic/copia-cli/pkg/cmd/release/create"
	deleteCmd "github.com/qubernetic/copia-cli/pkg/cmd/release/delete"
	listCmd "github.com/qubernetic/copia-cli/pkg/cmd/release/list"
	uploadCmd "github.com/qubernetic/copia-cli/pkg/cmd/release/upload"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
)

func NewCmdRelease(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release <command>",
		Short: "Manage releases",
		Long:  "Work with Copia repository releases.",
	}

	cmdutil.AddGroup(cmd, "General commands",
		listCmd.NewCmdList(f),
		createCmd.NewCmdCreate(f),
	)

	cmdutil.AddGroup(cmd, "Targeted commands",
		deleteCmd.NewCmdDelete(f),
		uploadCmd.NewCmdUpload(f),
	)

	return cmd
}
