package merge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

// MergeOptions holds all inputs for the pr merge command.
type MergeOptions struct {
	IO           *iostreams.IOStreams
	HTTPClient   *http.Client
	Host         string
	Token        string
	Owner        string
	Repo         string
	Number       int64
	Method       string
	DeleteBranch bool
}

type mergeRequest struct {
	Do                    string `json:"Do"`
	DeleteBranchAfterMerge bool   `json:"delete_branch_after_merge,omitempty"`
}

// NewCmdMerge creates the `copia pr merge` command.
func NewCmdMerge(f *cmdutil.Factory) *cobra.Command {
	opts := &MergeOptions{}

	cmd := &cobra.Command{
		Use:   "merge <number>",
		Short: "Merge a pull request",
		Example: `  copia pr merge 7
  copia pr merge 7 --squash
  copia pr merge 7 --rebase --delete-branch`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid PR number: %s", args[0])
			}
			opts.Number = num
			opts.IO = f.IOStreams

			host, token, err := f.ResolveAuth()
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token

			owner, repo, err := f.ResolveRepo()
			if err != nil {
				return err
			}
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}

			// Determine merge method from flags
			squash, _ := cmd.Flags().GetBool("squash")
			rebase, _ := cmd.Flags().GetBool("rebase")
			if squash {
				opts.Method = "squash"
			} else if rebase {
				opts.Method = "rebase"
			} else {
				opts.Method = "merge"
			}

			return mergeRun(opts)
		},
	}

	cmd.Flags().Bool("merge", false, "Merge commit (default)")
	cmd.Flags().Bool("squash", false, "Squash and merge")
	cmd.Flags().Bool("rebase", false, "Rebase and merge")
	cmd.Flags().BoolVar(&opts.DeleteBranch, "delete-branch", false, "Delete branch after merge")

	return cmd
}

func mergeRun(opts *MergeOptions) error {
	payload := mergeRequest{
		Do:                    opts.Method,
		DeleteBranchAfterMerge: opts.DeleteBranch,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/pulls/%d/merge",
		opts.Host, opts.Owner, opts.Repo, opts.Number)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to merge PR (HTTP %d)", resp.StatusCode)
	}

	_, _ = fmt.Fprintf(opts.IO.Out, "Merged PR #%d (%s)\n", opts.Number, opts.Method)
	return nil
}
