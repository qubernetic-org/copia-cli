package clone

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

// CloneOptions holds all inputs for the repo clone command.
type CloneOptions struct {
	IO    *iostreams.IOStreams
	Host  string
	Token string
	Repo  string
	Dir   string
}

// NewCmdClone creates the `copia repo clone` command.
func NewCmdClone(f *cmdutil.Factory) *cobra.Command {
	opts := &CloneOptions{}

	cmd := &cobra.Command{
		Use:   "clone <owner/repo | URL> [<directory>]",
		Short: "Clone a repository",
		Example: `  copia repo clone my-org/my-repo
  copia repo clone https://app.copia.io/my-org/my-repo.git`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.Repo = args[0]

			host, _, err := f.ResolveAuth()
			if err != nil {
				return err
			}
			opts.Host = host

			if len(args) > 1 {
				opts.Dir = args[1]
			}

			return cloneRun(opts)
		},
	}

	return cmd
}

func cloneRun(opts *CloneOptions) error {
	cloneURL := buildCloneURL(opts.Host, opts.Repo)

	args := []string{"clone", "--", cloneURL}
	if opts.Dir != "" {
		args = append(args, opts.Dir)
	}

	gitCmd := exec.Command("git", args...)
	gitCmd.Stdout = opts.IO.Out
	gitCmd.Stderr = opts.IO.ErrOut
	gitCmd.Stdin = os.Stdin

	return gitCmd.Run()
}

func buildCloneURL(host, repo string) string {
	if strings.HasPrefix(repo, "https://") || strings.HasPrefix(repo, "git@") {
		return repo
	}
	return fmt.Sprintf("https://%s/%s.git", host, repo)
}
