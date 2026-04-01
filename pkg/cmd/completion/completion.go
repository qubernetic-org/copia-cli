package completion

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

// NewCmdCompletion creates the `copia completion` command.
func NewCmdCompletion(ios *iostreams.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion <shell>",
		Short: "Generate shell completion scripts",
		Long: `Generate completion scripts for bash, zsh, fish, or powershell.

To load completions:

  # Bash
  source <(copia completion bash)

  # Zsh
  copia completion zsh > "${fpath[1]}/_copia"

  # Fish
  copia completion fish | source

  # PowerShell
  copia completion powershell | Out-String | Invoke-Expression`,
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletionV2(ios.Out, true)
			case "zsh":
				return cmd.Root().GenZshCompletion(ios.Out)
			case "fish":
				return cmd.Root().GenFishCompletion(ios.Out, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(ios.Out)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}

	return cmd
}
