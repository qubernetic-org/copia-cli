package org

import (
	"github.com/spf13/cobra"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/org/list"
	viewCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/org/view"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

func NewCmdOrg(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org <command>",
		Short: "Manage organizations",
		Long:  "Work with Copia organizations.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))
	cmd.AddCommand(viewCmd.NewCmdView(f))

	return cmd
}
