package copiacmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/internal/build"
	"github.com/qubernetic-org/copia-cli/internal/config"
	apiCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/api"
	authCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/auth"
	issueCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue"
	labelCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/label"
	prCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr"
	releaseCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/release"
	searchCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/search"
	repoCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/repo"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

// NewRootCmd creates the root `copia` command with all subcommands.
func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "copia <command> <subcommand> [flags]",
		Short:         "Copia CLI — source control for industrial automation",
		Long:          "Work with Copia repositories, issues, pull requests, and more from the command line.",
		Version:       build.Version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.SetVersionTemplate("copia version {{.Version}}\n")

	// Global flags
	cmd.PersistentFlags().StringVar(&f.Host, "host", "", "Target Copia host")
	cmd.PersistentFlags().StringVar(&f.Token, "token", "", "Authentication token")

	cmd.AddCommand(authCmd.NewCmdAuth(f))
	cmd.AddCommand(repoCmd.NewCmdRepo(f))
	cmd.AddCommand(issueCmd.NewCmdIssue(f))
	cmd.AddCommand(labelCmd.NewCmdLabel(f))
	cmd.AddCommand(prCmd.NewCmdPR(f))
	cmd.AddCommand(releaseCmd.NewCmdRelease(f))
	cmd.AddCommand(apiCmd.NewCmdApi(f))
	cmd.AddCommand(searchCmd.NewCmdSearch(f))

	return cmd
}

// Main is the entrypoint called from cmd/copia/main.go.
func Main() int {
	ios := iostreams.System()

	f := &cmdutil.Factory{
		IOStreams: ios,
		Config: func() (*config.Config, error) {
			return config.Load(config.DefaultPath())
		},
	}

	// Override from env if not set by flag
	if envToken := os.Getenv("COPIA_TOKEN"); envToken != "" {
		f.Token = envToken
	}
	if envHost := os.Getenv("COPIA_HOST"); envHost != "" {
		f.Host = envHost
	}

	rootCmd := NewRootCmd(f)

	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}
