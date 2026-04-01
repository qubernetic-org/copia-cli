package close

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

// CloseOptions holds all inputs for the issue close command.
type CloseOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	Number     int64
	Comment    string
}

// NewCmdClose creates the `copia issue close` command.
func NewCmdClose(f *cmdutil.Factory) *cobra.Command {
	opts := &CloseOptions{}

	cmd := &cobra.Command{
		Use:   "close <number>",
		Short: "Close an issue",
		Example: `  copia issue close 12
  copia issue close 12 --comment "Fixed in PR #7"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid issue number: %s", args[0])
			}
			opts.Number = num
			opts.IO = f.IOStreams

			host, token, err := f.ResolveAuth()
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token

			if f.BaseRepo == nil {
				return fmt.Errorf("could not determine repository. Run from inside a git repository")
			}
			owner, repo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return closeRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Comment, "comment", "c", "", "Add a comment before closing")

	return cmd
}

func closeRun(opts *CloseOptions) error {
	if opts.Comment != "" {
		commentPayload, _ := json.Marshal(map[string]string{"body": opts.Comment})
		commentURL := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues/%d/comments",
			opts.Host, opts.Owner, opts.Repo, opts.Number)

		req, err := http.NewRequest("POST", commentURL, bytes.NewReader(commentPayload))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "token "+opts.Token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := opts.HTTPClient.Do(req)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to add comment (HTTP %d)", resp.StatusCode)
		}
	}

	closePayload, _ := json.Marshal(map[string]string{"state": "closed"})
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues/%d",
		opts.Host, opts.Owner, opts.Repo, opts.Number)

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(closePayload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to close issue (HTTP %d)", resp.StatusCode)
	}

	_, _ = fmt.Fprintf(opts.IO.Out, "Closed issue #%d\n", opts.Number)
	return nil
}
