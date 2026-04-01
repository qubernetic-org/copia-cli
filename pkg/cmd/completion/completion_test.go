package completion

import (
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletion_Bash(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	root := &cobra.Command{Use: "copia"}
	root.AddCommand(NewCmdCompletion(ios))

	root.SetArgs([]string{"completion", "bash"})
	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "bash completion")
}

func TestCompletion_Zsh(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	root := &cobra.Command{Use: "copia"}
	root.AddCommand(NewCmdCompletion(ios))

	root.SetArgs([]string{"completion", "zsh"})
	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "zsh")
}

func TestCompletion_Fish(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	root := &cobra.Command{Use: "copia"}
	root.AddCommand(NewCmdCompletion(ios))

	root.SetArgs([]string{"completion", "fish"})
	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "fish")
}

func TestCompletion_PowerShell(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	root := &cobra.Command{Use: "copia"}
	root.AddCommand(NewCmdCompletion(ios))

	root.SetArgs([]string{"completion", "powershell"})
	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "copia")
}
