package notification

import (
	"github.com/spf13/cobra"
	listCmd "github.com/qubernetic/copia-cli/pkg/cmd/notification/list"
	readCmd "github.com/qubernetic/copia-cli/pkg/cmd/notification/read"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
)

func NewCmdNotification(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notification <command>",
		Short: "Manage notifications",
		Long:  "Work with Copia notifications.",
	}

	cmdutil.AddGroup(cmd, "General commands",
		listCmd.NewCmdList(f),
		readCmd.NewCmdRead(f),
	)

	return cmd
}
