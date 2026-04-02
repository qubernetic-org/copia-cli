package org

import (
	"github.com/spf13/cobra"
	listCmd "github.com/qubernetic/copia-cli/pkg/cmd/org/list"
	viewCmd "github.com/qubernetic/copia-cli/pkg/cmd/org/view"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
)

func NewCmdOrg(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org <command>",
		Short: "Manage organizations",
		Long:  "Work with Copia organizations.",
	}

	cmdutil.AddGroup(cmd, "General commands",
		listCmd.NewCmdList(f),
	)

	cmdutil.AddGroup(cmd, "Targeted commands",
		viewCmd.NewCmdView(f),
	)

	return cmd
}
