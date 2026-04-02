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
		Long:  "Clone a Copia repository locally. The repository can be specified as owner/repo or as a full URL.",
		Example: `  # Clone by owner/repo
  $ copia-cli repo clone my-org/my-repo

  # Clone by URL
  $ copia-cli repo clone https://app.copia.io/my-org/my-repo.git`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.Repo = args[0]

			host, token, err := f.ResolveAuth()
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token

			if len(args) > 1 {
				opts.Dir = args[1]
			}

			return CloneRun(opts)
		},
	}

	return cmd
}

func CloneRun(opts *CloneOptions) error {
	cloneURL := buildCloneURL(opts.Host, opts.Repo)

	// Use token in URL for authentication, then remove it from the remote.
	authURL := cloneURL
	if opts.Token != "" && !strings.HasPrefix(opts.Repo, "git@") {
		authURL = buildAuthCloneURL(opts.Host, opts.Token, opts.Repo)
	}

	args := []string{"clone", "--", authURL}
	dir := opts.Dir
	if dir != "" {
		args = append(args, dir)
	} else {
		// Derive dir name from repo for post-clone remote fix
		parts := strings.Split(strings.TrimSuffix(opts.Repo, ".git"), "/")
		dir = parts[len(parts)-1]
	}

	gitCmd := exec.Command("git", args...)
	gitCmd.Stdout = opts.IO.Out
	gitCmd.Stderr = opts.IO.ErrOut
	gitCmd.Stdin = os.Stdin
	gitCmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	if err := gitCmd.Run(); err != nil {
		return err
	}

	// Remove token from the stored remote URL
	if opts.Token != "" && authURL != cloneURL {
		setCmd := exec.Command("git", "-C", dir, "remote", "set-url", "origin", cloneURL)
		_ = setCmd.Run()
	}

	return nil
}

func buildCloneURL(host, repo string) string {
	if strings.HasPrefix(repo, "https://") || strings.HasPrefix(repo, "git@") {
		return repo
	}
	return fmt.Sprintf("https://%s/%s.git", host, repo)
}

func buildAuthCloneURL(host, token, repo string) string {
	if strings.HasPrefix(repo, "https://") {
		return strings.Replace(repo, "https://", fmt.Sprintf("https://token:%s@", token), 1)
	}
	return fmt.Sprintf("https://token:%s@%s/%s.git", token, host, repo)
}
