package completion

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

// NewCmdCompletion creates the `copia completion` command.
func NewCmdCompletion(ios *iostreams.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion <shell>",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for Copia CLI commands.

When installing Copia CLI through a package manager, it's possible that
no additional shell configuration is necessary to gain completion support. For
Homebrew, see <https://docs.brew.sh/Shell-Completion>

If you need to set up completions manually, follow the instructions below. The exact
config file locations might vary based on your system. Make sure to restart your
shell before testing whether completions are working.

### bash

First, ensure that you install ` + "`bash-completion`" + ` using your package manager.

After, add this to your ` + "`~/.bash_profile`" + `:

    eval "$(copia-cli completion bash)"

### zsh

Generate a ` + "`_copia-cli`" + ` completion script and put it somewhere in your ` + "`$fpath`" + `:

    copia-cli completion zsh > /usr/local/share/zsh/site-functions/_copia-cli

Ensure that the following is present in your ` + "`~/.zshrc`" + `:

    autoload -U compinit
    compinit -i

Zsh version 5.7 or later is recommended.

### fish

Generate a ` + "`copia-cli.fish`" + ` completion script:

    copia-cli completion fish > ~/.config/fish/completions/copia-cli.fish

### PowerShell

Open your profile script with:

    mkdir -Path (Split-Path -Parent $profile) -ErrorAction SilentlyContinue
    notepad $profile

Add the line and save the file:

    Invoke-Expression -Command $(copia-cli completion powershell | Out-String)`,
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
